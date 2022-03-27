package main

import (
	"log"
	"os"

	"github.com/aamilineni/go-products-review/constants"
	"github.com/aamilineni/go-products-review/db"
	"github.com/aamilineni/go-products-review/router"
	"github.com/aamilineni/go-products-review/server"
	"github.com/gin-gonic/gin"
)

func main() {

	// Check whether the username & password is supplied as
	username := os.Getenv(constants.AUTH_USERNAME)
	password := os.Getenv(constants.AUTH_PASSWORD)

	if username == "" {
		log.Fatal("basic auth username must be provided")
	}

	if password == "" {
		log.Fatal("basic auth password must be provided")
	}

	// Initialise Database
	db.InitDB()

	// Creates new default gin
	ginEngine := gin.Default()

	// Initialise Router
	router.InitialiseRouter(ginEngine)

	// Initialise Server with graceful shutdown
	server.InitialiseServer(ginEngine)
}
