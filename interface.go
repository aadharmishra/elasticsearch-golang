package elasticsearchgolang

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

type IDatabase interface {
	CreateDocument()
	GetDocument()
	UpdateDocument()
	DeleteDocument()
}

type esService struct {
	ctx *gin.Context
	es  *elasticsearch.Client
}
