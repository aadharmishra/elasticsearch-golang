package main

import (
	"fmt"
	"log"

	conn "esmodule/connection"
	handler "esmodule/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	//Initialise Elastic Client
	esClient, err := conn.InitElasticClient()

	if esClient == nil || err != nil || conn.ConnPool["client"] == nil {
		log.Fatalf("ES client init failed.")
		return
	}

	// Initialize Gin router
	router := gin.Default()

	// Define routes
	routerApi := router.Group("/author")
	{
		routerApi.POST("/create", handler.RequestHandler)
		routerApi.GET("/fetch", handler.RequestHandler)
		routerApi.PUT("/update", handler.RequestHandler)
		routerApi.PUT("/remove", handler.RequestHandler)
	}

	// Start server
	port := 8080
	router.Run(fmt.Sprintf(":%d", port))
}
