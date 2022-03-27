package products

import (
	"github.com/gin-gonic/gin"
)

var RegisterProductsReviewRoutes = func(r *gin.RouterGroup) {

	productsHandler := NewProductsHandler()
	r.GET("/products", productsHandler.Get)
}
