package repository

import (
	"context"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/aamilineni/go-products-review/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// creates new productsRepo interface
type ProductsRepo interface {
	GetAll()
	Add(product *models.ProductRequestModel) (string, error)
	AddReview(product *models.ProductReview)
}

type productsRepository struct {
	db *mongo.Database
}

func NewProductsRepository() ProductsRepo {
	return &productsRepository{
		db: db.GetMongoClient().Database(db.DB),
	}
}

// GetAll the products
func (me *productsRepository) GetAll() {

}

// Add the new product
// returns the insertedID & error
func (me *productsRepository) Add(product *models.ProductRequestModel) (string, error) {

	//Create a handle to the respective collection in the database.
	collection := me.db.Collection(db.PRODUCTSCOLLECTION)
	//Perform InsertOne operation & validate against the error.
	result, err := collection.InsertOne(context.TODO(), product)
	if err != nil {
		return "", err
	}

	// Return success without any error.
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	// Return err
	return "", nil
}

// Add the new product review
func (me *productsRepository) AddReview(product *models.ProductReview) {

}
