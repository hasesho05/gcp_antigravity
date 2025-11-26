package firestore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context) *firestore.Client {
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	conf := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT_ID")}
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}
	return client
}
