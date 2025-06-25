package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var Collection *mongo.Collection
var UserCollection *mongo.Collection

func ConnectDB() {
	// Load environment variables (optional for Docker)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file:", err)
		log.Println("Continuing with system environment variables...")
	}

	// Retrieve the connection string
	connectionString := os.Getenv("ConnectionString")
	if connectionString == "" {
		log.Fatal("Missing ConnectionString in .env")
	}

	// Set client options
	clientOption := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	// Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}

	fmt.Println("MongoDB connection success")
	dbName := getEnvOrDefault("DBName", "anime")
	colName := getEnvOrDefault("CollectionName", "watchlist")
	userColName := getEnvOrDefault("UserCollectionName", "users")
	Collection = client.Database(dbName).Collection(colName)
	UserCollection = client.Database(dbName).Collection(userColName)
}
