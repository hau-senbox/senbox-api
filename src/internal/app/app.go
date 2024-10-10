package app

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"os"
	"os/signal"
	"sen-global-api/config"
	"sen-global-api/docs"
	"sen-global-api/internal/database"
	"sen-global-api/internal/middleware"
	"sen-global-api/internal/router"
	"sen-global-api/pkg/common"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/mysql"
	"sen-global-api/pkg/sheet"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(appConfig *config.AppConfig, fcm *firebase.App) error {
	monitor.SendMessageViaTelegram("Server is starting...")
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

	monitor.SendMessageViaTelegram("Server is up and running...")

	select {
	case s := <-interrupt:
		monitor.SendMessageViaTelegram("Server is interrupting...", s.String())
		log.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		monitor.SendMessageViaTelegram("Server is shutting down...", err.Error())
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		monitor.SendMessageViaTelegram("Server is shutting down...", err.Error())
	}

	return err
}
