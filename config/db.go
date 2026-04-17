package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// DB is our global handle to the database client.
var DB *mongo.Client

func Connect() {
	// 1. Configure options with our connection string.
	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	// 2. Connect to the database.
	// In v2, Connect does not require a context directly.
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// 3. Ping the database with a 5-second timeout to confirm it's reachable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB v2!")
	DB = client
}
