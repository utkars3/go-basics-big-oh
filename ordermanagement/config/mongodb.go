package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var mongoClient *mongo.Client

func ConnectToMongoDB() error {
	// Load environment variables

	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get MongoDB URI from .env
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not found in environment variables")
	}

	// Define MongoDB connection options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB (newer method)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to ensure connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	// log.Println("✅ Connected to MongoDB!")

	// Set global MongoDB client and database
	mongoClient = client
	MongoDB = client.Database("ordermanagement")
	log.Println("MongoDB instance:", MongoDB)

	log.Println("✅ Connected to MongoDB!")
	return nil
}

// DisconnectMongoDB closes the MongoDB connection
func DisconnectMongoDB() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		} else {
			log.Println("🛑 MongoDB connection closed.")
		}
	}
}
