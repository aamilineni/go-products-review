package products

import (
	"net/http"

	"github.com/aamilineni/go-products-review/api/models"
	"github.com/aamilineni/go-products-review/api/repository"
	"github.com/gin-gonic/gin"
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
		ctx.AbortWithError(http.StatusBadRequest, models.AppError{
			Status:       http.StatusBadRequest,
			ErrorMessage: "required fields cannot be empty",
		})
		return
	}

	var err error
	// Set the product id created in the db
	productRequestModel.ID, err = me.productRepo.Add(&productRequestModel)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, models.AppError{
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

}

// Add new review to the product
func (me *productsHandler) AddReview(ctx *gin.Context) {

}
