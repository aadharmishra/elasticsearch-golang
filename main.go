package elasticsearchgolang

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

var alreadyPushedData bool

var connPool map[string]*elasticsearch.Client

var isTypedApiFlowEnabled bool = true

func main() {

	// Initialize Gin router
	router := gin.Default()

	// Define routes
	routerApi := router.Group("/author")
	{
		routerApi.POST("/create", RequestHandler)
		routerApi.GET("/fetch", RequestHandler)
		routerApi.PUT("/update", RequestHandler)
		routerApi.PUT("/remove", RequestHandler)
	}

	//Initialise Elastic Client
	esClient, err := InitElasticClient()

	if esClient == nil || err != nil {
		log.Fatalf("Client init failed.")
		return
	}

	//global connection pool
	connPool = make(map[string]*elasticsearch.Client)
	connPool["client"] = esClient

	// Start server
	port := 8080
	router.Run(fmt.Sprintf(":%d", port))
}
