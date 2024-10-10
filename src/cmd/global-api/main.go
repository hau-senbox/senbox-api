package main

import (
	"log"
	"os"
	config2 "sen-global-api/config"
	"sen-global-api/internal/app"
	"sen-global-api/pkg/logger"
	"sen-global-api/pkg/messaging"

	"github.com/ilyakaznacheev/cleanenv"
)

// @title           Swagger Sen Master API
// @version         1.0
// @BasePath  /api/v1

func main() {

	// Read config from args and init the app config
	configFilePath := os.Args[1]
	appConfig := &config2.AppConfig{}

	err := cleanenv.ReadConfig(configFilePath, appConfig)
	if err != nil {
		log.Panic(err)
	}

	//Config Logger
	err = logger.InitLogger(appConfig)
	if err != nil {
		log.Panic(err)
	}

	fcm, err := messaging.NewFirebaseApp(*appConfig)
	if err != nil {
		log.Panic(err)
	}

	// Run the app
	err = app.Run(appConfig, fcm)
	if err != nil {
		log.Panic(err)
	}
}
