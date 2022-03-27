package repository

import (
	"context"
	"math"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/aamilineni/go-products-review/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// creates new productsRepo interface
type ProductsRepo interface {
	GetBy(productName string, pageNumber int64, limit int64) (*models.ProductPaginationResponse, error)
	Add(product *models.Product) (string, error)
	AddReview(productID string, review *models.ProductReview) (string, error)
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
func (me *productsRepository) GetBy(productName string, pageNumber int64, limit int64) (*models.ProductPaginationResponse, error) {
	//Create a handle to the respective collection in the database.
	collection := me.db.Collection(db.PRODUCTSCOLLECTION)

	ctx := context.TODO()
	// query to find all the products data

	filter := bson.M{}
	if productName != "" {
		filter = bson.M{
			"name": bson.M{
				"$regex": primitive.Regex{
					Pattern: productName,
					Options: "i",
				},
			},
		}
	}

	// find the count of all the documents in a collection
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// set pagination options
	findOptions := &options.FindOptions{}
	findOptions.SetSkip((pageNumber - 1) * limit)
	findOptions.SetLimit(limit)

	// fetch from the collection based on the pagination & filtering options
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var results []models.Product
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return &models.ProductPaginationResponse{
		Total:    total,
		Page:     pageNumber,
		Limit:    limit,
		LastPage: int64(math.Ceil(float64(total) / float64(limit))),
		Products: results,
	}, nil
}

// Add the new product
// returns the insertedID & error
func (me *productsRepository) Add(product *models.Product) (string, error) {

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
func (me *productsRepository) AddReview(productID string, review *models.ProductReview) (string, error) {

	//Create a handle to the respective collection in the database.
	collection := me.db.Collection(db.PRODUCTSCOLLECTION)
	//Perform InsertOne operation & validate against the error.
	id, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return "", err
	}

	update := bson.M{
		"$push": bson.M{
			"reviews": bson.M{"$each": []models.ProductReview{{
				ID:                primitive.NewObjectID(),
				ReviewerName:      review.ReviewerName,
				ReviewDescription: review.ReviewDescription,
				ReviewRating:      review.ReviewRating,
			}}},
		},
	}
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		update,
	)
	if err != nil {
		return "", err
	}

	// Return success without any error.
	if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	// Return err
	return "", nil
}
