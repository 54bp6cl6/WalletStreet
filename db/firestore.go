package db

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

var (
	client *firestore.Client
	ctx    context.Context
)

func Connect() {
	var err error
	ctx = context.Background()
	client, err = firestore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
}

func Close() {
	ctx.Done()
	client.Close()
}
