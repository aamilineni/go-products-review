package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

//I have used below constants just to hold required database config's.
const (
	CONNECTIONSTRING   = "mongodb://mongodb"
	DB                 = "db_products"
	PRODUCTSCOLLECTION = "products"
)

//GetMongoClient - Return mongodb connection to work with
func GetMongoClient() *mongo.Client {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatalf("error while creating the mongo connection :: %+v", err)

			return
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatalf("error while pinging the mongo db :: %+v", err)

			return
		}
		clientInstance = client
	})
	return clientInstance
}
