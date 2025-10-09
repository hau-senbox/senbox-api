package main

import (
	"context"
	"log"
	"os"
	config2 "sen-global-api/config"
	"sen-global-api/internal/app"
	"sen-global-api/pkg/logger"
	"sen-global-api/pkg/messaging"

	"cloud.google.com/go/firestore"
	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/api/option"
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

	// firestore firebase setup

	ctx := context.Background()

	sa := option.WithCredentialsFile("credentials/senboxapp-firebase-adminsdk.json") // đường dẫn tới file key
	client, err := firestore.NewClient(ctx, "senboxapp-a1ad0", sa)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	if err != nil {
		log.Fatalf("Failed to add document: %v", err)
	}
	log.Println("Document added")

	// Run the app
	err = app.Run(appConfig, fcm)
	if err != nil {
		log.Panic(err)
	}
}
