package handlers

import (
	"strings"

	conn "esmodule/connection"
	"esmodule/interfaces"

	"esmodule/clients"

	"github.com/gin-gonic/gin"
)

func RequestHandler(ctx *gin.Context) {
	url := ctx.Request.URL.String()
	method := ctx.Request.Method

	esClient := conn.ConnPool["client"]

	if esClient == nil {
		return
	}

	var service interfaces.IDatabase

	service = &clients.EsService{
		Ctx: ctx,
		Es:  esClient,
	}

	if strings.Contains(url, "/create") && method == "POST" {
		service.CreateDocument()
	}
	if strings.Contains(url, "/fetch") && method == "GET" {
		service.GetDocument()
	}
	if strings.Contains(url, "/update") && method == "PUT" {
		service.UpdateDocument()
	}
	if strings.Contains(url, "/remove") && method == "PUT" {
		service.DeleteDocument()
	}
}
