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

const dbName = "anime"
const colName = "watchlist"

var Collection *mongo.Collection

func ConnectDB() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
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
	Collection = client.Database(dbName).Collection(colName)
}
