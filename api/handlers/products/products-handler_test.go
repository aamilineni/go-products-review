package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockProductRepo struct {
	gomock.Controller
}

type GetBy func(productName string, pageNumber int64, limit int64) (*models.ProductPaginationResponse, error)
type Add func(product *models.Product) (string, error)
type AddReview func(productID string, review *models.ProductReview) (string, error)

var getBy GetBy
var add Add
var addReview AddReview

func (me *mockProductRepo) GetBy(productName string, pageNumber int64, limit int64) (*models.ProductPaginationResponse, error) {
	return getBy(productName, pageNumber, limit)
}

func (me *mockProductRepo) Add(product *models.Product) (string, error) {
	return add(product)
}

func (me *mockProductRepo) AddReview(productID string, review *models.ProductReview) (string, error) {
	return addReview(productID, review)
}

func Test_Add_For_Success_With_Status_Code_201(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRepo := &mockProductRepo{}
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	product := models.Product{
		Name:        "Product 1",
		Description: "description",
	}
	add = func(product *models.Product) (string, error) {
		return "id", nil
	}
	requestBytes, err := json.Marshal(product)
	assert.NoError(t, err)
	ctx.Request, err = http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(requestBytes))
	assert.NoError(t, err)
	ctx.Request.Header.Add("Content-Type", "application/json")

	r.POST("/product", NewProductsHandler(mockRepo).Add)
	r.ServeHTTP(w, ctx.Request)

	assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
}

func Test_GetBy_For_Success_With_Products_Count_100(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRepo := &mockProductRepo{}
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	getBy = func(productName string, pageNumber, limit int64) (*models.ProductPaginationResponse, error) {
		return &models.ProductPaginationResponse{
			Total:    100,
			Page:     1,
			Limit:    10,
			LastPage: 10,
			Products: getMockProducts(100),
		}, nil
	}
	var err error
	ctx.Request, err = http.NewRequest(http.MethodGet, "/products", nil)
	assert.NoError(t, err)
	ctx.Request.Header.Add("Content-Type", "application/json")

	r.GET("/products", NewProductsHandler(mockRepo).Get)
	r.ServeHTTP(w, ctx.Request)

	var response models.ProductPaginationResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, ctx.Writer.Status(), http.StatusOK)
	assert.Equal(t, len(response.Products), 100)
}

func Test_AddReview_For_Success(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRepo := &mockProductRepo{}
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	addReview = func(productID string, review *models.ProductReview) (string, error) {
		return "id", nil
	}

	review := &models.ProductReviewRequestModel{
		ReviewerName:      "Anil",
		ReviewDescription: "Desc",
		ReviewRating:      5,
	}
	requestBytes, err := json.Marshal(review)
	assert.NoError(t, err)
	ctx.Request, err = http.NewRequest(http.MethodPost, "/product/productid/review", bytes.NewReader(requestBytes))
	assert.NoError(t, err)
	ctx.Request.Header.Add("Content-Type", "application/json")

	r.POST("/product/:id/review", NewProductsHandler(mockRepo).AddReview)
	r.ServeHTTP(w, ctx.Request)

	assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
}

func Test_AddReview_For_Invalid_Params(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRepo := &mockProductRepo{}
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	addReview = func(productID string, review *models.ProductReview) (string, error) {
		return "id", nil
	}

	review := &models.ProductReviewRequestModel{
		ReviewerName:      "Anil",
		ReviewDescription: "Desc",
		ReviewRating:      10,
	}
	requestBytes, err := json.Marshal(review)
	assert.NoError(t, err)
	ctx.Request, err = http.NewRequest(http.MethodPost, "/product/productid/review", bytes.NewReader(requestBytes))
	assert.NoError(t, err)
	ctx.Request.Header.Add("Content-Type", "application/json")

	r.POST("/product/:id/review", NewProductsHandler(mockRepo).AddReview)
	r.ServeHTTP(w, ctx.Request)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func Test_AddReview_For_Invalid_Product(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRepo := &mockProductRepo{}
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	addReview = func(productID string, review *models.ProductReview) (string, error) {
		return "", fmt.Errorf("invalid product id")
	}
	review := &models.ProductReviewRequestModel{
		ReviewerName:      "Anil",
		ReviewDescription: "Desc",
		ReviewRating:      5,
	}
	requestBytes, err := json.Marshal(review)
	assert.NoError(t, err)
	ctx.Request, err = http.NewRequest(http.MethodPost, "/product/productid/review", bytes.NewReader(requestBytes))
	assert.NoError(t, err)
	ctx.Request.Header.Add("Content-Type", "application/json")

	r.POST("/product/:id/review", NewProductsHandler(mockRepo).AddReview)
	r.ServeHTTP(w, ctx.Request)

	assert.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)
}

func getMockProducts(count int) []models.Product {

	faker := faker.New()
	products := []models.Product{}
	for i := 0; i < count; i++ {
		products = append(products, models.Product{
			ID:          primitive.NewObjectID(),
			Name:        faker.Company().Name(),
			Description: faker.Lorem().Text(50),
			Reviews: []models.ProductReview{
				{
					ID:                primitive.NewObjectID(),
					ReviewerName:      faker.Beer().Name(),
					ReviewDescription: faker.Lorem().Text(50),
					ReviewRating:      faker.IntBetween(0, 5),
				},
				{
					ID:                primitive.NewObjectID(),
					ReviewerName:      faker.Beer().Name(),
					ReviewDescription: faker.Lorem().Text(50),
					ReviewRating:      faker.IntBetween(0, 5),
				},
			},
		})
	}
	return products
}
