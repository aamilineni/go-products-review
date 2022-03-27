package router

import (
	"net/http"

	"github.com/aamilineni/go-products-review/api/handlers/products"
	"github.com/aamilineni/go-products-review/middleware"
	"github.com/gin-gonic/gin"
)

func InitialiseRouter(engine *gin.Engine) {

	// create the versioning of API
	apiV1 := engine.Group("/api/v1")

	// Set the auth middleware
	apiV1.Use(middleware.AuthMiddleware())

	// health check API
	apiV1.GET("/healthcheck", healthcheck)

	// Register Product handlers
	products.RegisterProductsReviewRoutes(apiV1)

}

func healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
