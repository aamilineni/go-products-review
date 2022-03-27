package products

import "github.com/gin-gonic/gin"

type productsHandler struct {
}

// NewProductsHandler creates new product handler
func NewProductsHandler() *productsHandler {
	return &productsHandler{}
}

// Get all the products
func (me *productsHandler) Get(ctx *gin.Context) {

}
