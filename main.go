package main

import (
	"github.com/aamilineni/go-products-review/router"
	"github.com/aamilineni/go-products-review/server"
	"github.com/gin-gonic/gin"
)

func main() {

	// Creates new default gin
	ginEngine := gin.Default()

	// Initialise Router
	router.InitialiseRouter(ginEngine)

	// Initialise Server with graceful shutdown
	server.InitialiseServer(ginEngine)
}
