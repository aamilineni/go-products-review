package products

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/aamilineni/go-products-review/api/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productsHandler struct {
	productRepo repository.ProductsRepo
}

// NewProductsHandler creates new product handler
func NewProductsHandler(repo repository.ProductsRepo) *productsHandler {
	return &productsHandler{
		productRepo: repo,
	}
}

// Add new product
func (me *productsHandler) Add(ctx *gin.Context) {

	// bind the request body
	var productRequestModel models.ProductRequestModel
	if err := ctx.BindJSON(&productRequestModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.AppError{
			Status:       http.StatusBadRequest,
			ErrorMessage: "required fields cannot be empty",
		})
		return
	}

	var err error
	// Set the product id created in the db
	productRequestModel.ID, err = me.productRepo.Add(&models.Product{
		ID:           primitive.NewObjectID(),
		Name:         productRequestModel.Name,
		Description:  productRequestModel.Description,
		ThumbnailURL: productRequestModel.ThumbnailURL,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.AppError{
			Status:       http.StatusInternalServerError,
			ErrorMessage: "failed to insert the product",
		})
		return
	}

	// response with status created and the product requested to create
	ctx.JSON(http.StatusCreated, productRequestModel)

}

// Get all the products
func (me *productsHandler) Get(ctx *gin.Context) {

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "5")
	productName := ctx.Query("name")

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.AppError{
			Status:       http.StatusInternalServerError,
			ErrorMessage: "page number in the query should be valid integer",
		})
		return
	}
	limitCount, err := strconv.Atoi(limit)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.AppError{
			Status:       http.StatusInternalServerError,
			ErrorMessage: "limit in the query should be valid integer",
		})
		return
	}

	// get all the products
	results, err := me.productRepo.GetBy(productName, int64(pageNumber), int64(limitCount))
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.AppError{
			Status:       http.StatusInternalServerError,
			ErrorMessage: "failed to get the products",
		})
		return
	}

	// response with status created and the product requested to create
	ctx.JSON(http.StatusOK, results)
}

// Add new review to the product
func (me *productsHandler) AddReview(ctx *gin.Context) {

	// Get the product id from the URL path
	productID := ctx.Param("id")

	// bind the request body
	var reviewRequestModel models.ProductReviewRequestModel
	if err := ctx.BindJSON(&reviewRequestModel); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.AppError{
			Status:       http.StatusBadRequest,
			ErrorMessage: "required fields cannot be empty",
		})
		return
	}

	var err error
	reviewUpsert := &models.ProductReview{
		ID:                primitive.NewObjectID(),
		ReviewerName:      reviewRequestModel.ReviewerName,
		ReviewDescription: reviewRequestModel.ReviewDescription,
		ReviewRating:      reviewRequestModel.ReviewRating,
	}
	// Set the product id created in the db
	_, err = me.productRepo.AddReview(productID, reviewUpsert)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.AppError{
			Status:       http.StatusInternalServerError,
			ErrorMessage: "failed to add the product review",
		})
		return
	}

	// response with status created and the product requested to create
	ctx.JSON(http.StatusCreated, reviewUpsert)

}
