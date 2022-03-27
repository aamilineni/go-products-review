package router

import (
	"net/http"

	"github.com/aamilineni/go-products-review/api/handlers/products"
	"github.com/aamilineni/go-products-review/api/handlers/productsreview"
	"github.com/gin-gonic/gin"
)

func InitialiseRouter(engine *gin.Engine) {

	// create the versioning of API
	apiV1 := engine.Group("/api/v1")

	// health check API
	apiV1.GET("/healthcheck", healthcheck)

	// Register Product handlers
	products.RegisterProductsReviewRoutes(apiV1)

	// Register Product Review handlers
	productsreview.RegisterProductsReviewRoutes(apiV1)

}

func healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}