package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func InitDB() {

	ctx := context.TODO()

	// Create an Index on `name` in products collection
	collection := GetMongoClient().Database(DB).Collection(PRODUCTSCOLLECTION)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"name": 1},
		Options: &options.IndexOptions{
			Collation: &options.Collation{
				Locale:   "en",
				Strength: 2,
			},
		},
	}, opts)
	log.Println("Successfully create the index")

	// Seed products
	total, err := GetMongoClient().Database(DB).Collection(PRODUCTSCOLLECTION).CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("error while getting count of products collection")
		return
	}
	if total == 0 {
		seedProducts(ctx)
	}

}

func seedProducts(ctx context.Context) {

	// Seed the products
	faker := faker.New()
	products := []interface{}{}
	for i := 0; i < 50; i++ {
		product := models.Product{
			ID:           primitive.NewObjectID(),
			Name:         faker.Company().Name(),
			Description:  faker.Lorem().Paragraph(2),
			ThumbnailURL: faker.Internet().URL(),
			Reviews:      []models.ProductReview{},
		}
		products = append(products, product)
	}

	result, err := GetMongoClient().Database(DB).Collection(PRODUCTSCOLLECTION).InsertMany(ctx, products)
	if err != nil {
		log.Fatal("error while seeding into products collection")
		return
	}
	log.Println("Total Seeded Products are :: ", len(result.InsertedIDs))
}
