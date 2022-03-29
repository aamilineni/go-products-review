package repository

import (
	"context"
	"math"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/aamilineni/go-products-review/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		filter = bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: productName, Options: "i"}}}
	}

	// find the count of all the documents in a collection
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	var basePipeline = mongo.Pipeline{}

	if productName != "" {
		subpipeline := bson.D{{"$match", bson.D{{"$and", bson.A{bson.D{{"name", bson.M{"$regex": primitive.Regex{Pattern: productName, Options: "i"}}}}}}}}}
		basePipeline = append(basePipeline, subpipeline)
	}

	// lookup for join
	lookup := bson.D{
		{"$lookup", bson.D{
			{"from", "reviews"},
			{"localField", "_id"},
			{"foreignField", "productID"},
			{"as", "averagereviews"},
		}},
	}
	basePipeline = append(basePipeline, lookup)

	// unwind
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$averagereviews"}, {"preserveNullAndEmptyArrays", true}}}}
	basePipeline = append(basePipeline, unwindStage)

	// group
	group := bson.D{
		{
			"$group", bson.D{
				{"_id", "$_id"},
				{"name", bson.D{{"$first", "$name"}}},
				{"description", bson.D{{"$first", "$description"}}},
				{"thumbnail_url", bson.D{{"$first", "$thumbnail_url"}}},
				{"average_rating", bson.D{{"$avg", "$averagereviews.rating"}}},
			},
		},
	}
	basePipeline = append(basePipeline, group)

	// sort on id
	sort := bson.D{
		{
			"$sort", bson.M{"_id": 1},
		},
	}
	basePipeline = append(basePipeline, sort)

	// for pagination
	jumpPagePipeline := bson.D{{"$skip", (pageNumber - 1) * limit}}
	limitPipeline := bson.D{{"$limit", limit}}
	basePipeline = append(basePipeline, jumpPagePipeline)
	basePipeline = append(basePipeline, limitPipeline)

	// fetch from the collection based on the pagination & filtering options
	cursor, err := collection.Aggregate(ctx, basePipeline)
	if err != nil {
		return nil, err
	}

	var results []models.ProductResponseModel
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
	collection := me.db.Collection(db.REVIEWCOLLECTION)
	//Perform InsertOne operation & validate against the error.
	id, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return "", err
	}

	review.ProductID = id
	result, err := collection.InsertOne(context.TODO(), review)
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
