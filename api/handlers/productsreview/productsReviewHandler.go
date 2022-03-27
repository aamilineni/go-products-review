package productsreview

import "github.com/gin-gonic/gin"

type productsReviewHandler struct {
}

// NewProductsReviewHandler creates new product review handler
func NewProductsReviewHandler() *productsReviewHandler {
	return &productsReviewHandler{}
}

// Add new product review
func (me *productsReviewHandler) Add(ctx *gin.Context) {

}
