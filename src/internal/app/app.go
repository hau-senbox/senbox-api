package app

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sen-global-api/config"
	"sen-global-api/docs"
	"sen-global-api/internal/database"
	"sen-global-api/internal/middleware"
	"sen-global-api/internal/router"
	"sen-global-api/pkg/common"
	"sen-global-api/pkg/mysql"
	"sen-global-api/pkg/sheet"
	"strconv"
	"syscall"
	"time"

	"sen-global-api/internal/domain/usecase"

	firebase "firebase.google.com/go/v4"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	serviceName = "go-main-service"
	ttl         = time.Second * 8
	checkId     = "go-main-service-health-check"
)

var serviceId = fmt.Sprintf("%s-%d", serviceName, rand.Intn(100))

func Run(appConfig *config.AppConfig, fcm *firebase.App) error {
	if appConfig.Config.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//Establish connection to database
	dbConn, err := mysql.Establish(*appConfig)
	if err != nil {
		log.Fatal("Could not connect to database ", err)
	}
	//err = migrations.MigrateDevices(dbConn)
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = database.Seed(dbConn, appConfig.Config, "/internal/database/seed.sql")
	if err != nil {
		log.Fatal(err)
	}

	//Establish connection to google sheet
	ctx := context.Background()
	userSpreadsheet, err := sheet.NewUserSpreadsheet(*appConfig, ctx)

	if err != nil {
		log.Fatal(err)
	}

	uploaderSpreadsheet, err := sheet.NewUploaderSpreadsheet(*appConfig, ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Initial server
	handler := gin.New()
	handler.Use(gin.CustomRecovery(middleware.RecoveryHandler), middleware.CORS())
	router.Route(handler, dbConn, userSpreadsheet, uploaderSpreadsheet, *appConfig, fcm)

	docs.SwaggerInfo.BasePath = "/"
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	httpServer := common.NewServer(handler, common.Port(appConfig.Config.HTTP.Port))

	log.Debug(fmt.Sprintf("Starting HTTP server on port %s", appConfig.Config.HTTP.Port))
	log.Info(fmt.Sprintf("Starting HTTP server on port %s", appConfig.Config.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	consulHost := appConfig.Config.Consul.Host
	if consulHost == "" {
		// Fallback to localhost if the host is not set in the config
		consulHost = "localhost"
	}

	// Consul client setup
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%s", consulHost, appConfig.Config.Consul.Port), // Consul server address
		HttpClient: &http.Client{
			Timeout: 30 * time.Second, // Increase timeout to 30 seconds
		},
	})

	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	setupConsul(client, consulHost, appConfig)
	go updateHealthCheck(client)
	usecase.ConsulClient = client
	handler.GET("/health", healthCheck)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
		deregisterConsul(client)
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		deregisterConsul(client)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		deregisterConsul(client)
	}

	return err
}

func healthCheck(c *gin.Context) {
	c.Status(200)
}

func updateHealthCheck(client *api.Client) {
	ticker := time.NewTicker(time.Second * 5)

	for {
		err := client.Agent().UpdateTTL(checkId, "online", api.HealthPassing)
		// _, _, err := client.Agent().AgentHealthServiceByID(serviceId)
		if err != nil {
			// log.Fatalf("Failed to update health check: %v", err)
			log.Fatalf("Failed to check AgentHealthService: %v", err)
		}
		<-ticker.C
	}
}

func setupConsul(client *api.Client, consulHost string, appConfig *config.AppConfig) {
	hostname := appConfig.Config.Registry.Host
	// hostname, _ := os.Hostname()
	port, _ := strconv.Atoi(appConfig.Config.HTTP.Port)

	// healthCheckHost := "localhost"
	// if appConfig.Config.Consul.Host != "localhost" {
	// 	// Fallback to localhost if the host is not set in the config
	// 	healthCheckHost = serviceName
	// }

	// Health check (optional but recommended)
	check := &api.AgentServiceCheck{
		// HTTP:     fmt.Sprintf("http://%s:%v/health", healthCheckHost, port), // Health check endpoint
		// Interval: "10s",                                                     // Interval for health check
		// Timeout:  "30s",
		DeregisterCriticalServiceAfter: ttl.String(),
		TTL:                            ttl.String(),
		CheckID:                        checkId,
	}

	// Service registration
	registration := &api.AgentServiceRegistration{
		ID:      serviceId,   // Unique service RoleId
		Name:    serviceName, // Service name
		Port:    port,        // Service port
		Address: hostname,    // Service address
		Tags:    []string{"go", "main", "user-service"},
		Check:   check,
	}

	query := map[string]any{
		"type":        "service",
		"service":     "mycluster",
		"passingonly": true,
	}

	plan, err := watch.Parse(query)
	if err != nil {
		log.Fatalf("Failed to watch for changes: %v", err)
	}

	plan.HybridHandler = func(index watch.BlockingParamVal, result any) {
		switch msg := result.(type) {
		case []*api.ServiceEntry:
			for _, entry := range msg {
				println("new member joined: ", entry.Service)
			}
		}
	}

	go func() {
		_ = plan.RunWithConfig(fmt.Sprintf("%s:%s", consulHost, appConfig.Config.Consul.Port), api.DefaultConfig())
	}()

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Panic(err)
		log.Printf("Failed to register service: %s:%v ", hostname, port)
		log.Fatalf("Failed to register health check: %v", err)
	}

	log.Printf("successfully register service: %s:%v", hostname, port)
}

func deregisterConsul(client *api.Client) {
	// Deregister service
	err := client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		log.Fatalf("Failed to deregister service: %v", err)
	}
}
