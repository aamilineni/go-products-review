package products

import (
	"github.com/aamilineni/go-products-review/api/repository"
	"github.com/gin-gonic/gin"
)

var RegisterProductsReviewRoutes = func(r *gin.RouterGroup) {

	productsHandler := NewProductsHandler(repository.NewProductsRepository())

	r.GET("/products", productsHandler.Get)
	r.POST("/product", productsHandler.Add)
	r.POST("/product/:id/review", productsHandler.AddReview)
}
