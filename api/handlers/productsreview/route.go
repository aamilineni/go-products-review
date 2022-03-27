package productsreview

import (
	"github.com/gin-gonic/gin"
)

var RegisterProductsReviewRoutes = func(r *gin.RouterGroup) {

	productsHandler := NewProductsReviewHandler()
	r.GET("/products", productsHandler.Add)
}
