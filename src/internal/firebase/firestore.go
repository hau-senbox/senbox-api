package firebase

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	once            sync.Once
)

func InitFirestoreClient() *firestore.Client {
	once.Do(func() {
		ctx := context.Background()
		opt := option.WithCredentialsFile("credentials/senboxapp-firebase-adminsdk.json")
		client, err := firestore.NewClient(ctx, "senboxapp-a1ad0", opt)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore: %v", err)
		}
		firestoreClient = client
	})
	return firestoreClient
}
