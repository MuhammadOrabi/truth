package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
)

// Database connection
type Database struct {
	Client  *mongo.Client
	DB      *mongo.Database
	Sources *mongo.Collection
	Users   *mongo.Collection
}

var mongoOnce sync.Once
var clientInstanceError error
var db Database

// Load ...
func Load() *Database {
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		db.Client = client
	})

	if clientInstanceError != nil {
		log.Fatal(clientInstanceError)
	}

	db.DB = db.Client.Database(os.Getenv("MONGO_DB"))
	db.Sources = db.DB.Collection("sources")
	db.Users = db.DB.Collection("users")

	return &db
}
