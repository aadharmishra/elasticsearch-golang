package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"esmodule/models"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

var isTypedApiFlowEnabled bool = true

type EsService struct {
	Ctx *gin.Context
	Es  *elasticsearch.Client
}

// Create document
func (service *EsService) CreateDocument() {

	var authorReq *models.Author
	var res *esapi.Response
	var err error
	var success bool

	ctx := service.Ctx
	es := service.Es

	if ctx == nil || es == nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = ctx.ShouldBindJSON(&authorReq)
	if err != nil || authorReq == nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	exists, err := DocumentExists(ctx, es, authorReq)
	if err != nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	if exists {
		fmt.Printf("Document with ID %s already exists. Skipping insertion.\n", authorReq.Details.DocumentID)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	doc := map[string]interface{}{
		"index": map[string]interface{}{
			"_id": authorReq.Details.DocumentID,
		},
		"name":         authorReq.Name,
		"author":       authorReq.Author,
		"release_date": authorReq.ReleaseDate,
		"page_count":   authorReq.PageCount,
	}

	jsonUpdateDoc, err := json.Marshal(doc)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	if !isTypedApiFlowEnabled {
		res, err = es.Bulk(
			bytes.NewReader(jsonUpdateDoc),
			es.Bulk.WithIndex(authorReq.Details.Index),
		)
	} else {
		req := esapi.IndexRequest{
			Index:      authorReq.Details.Index,
			DocumentID: authorReq.Details.DocumentID,
			Body:       bytes.NewReader(jsonUpdateDoc),
		}

		res, err = req.Do(context.Background(), es)
	}

	defer res.Body.Close()

	success, err = ValidateResponse(res, err)
	if !success || err != nil {
		ctx.JSON(res.StatusCode, err)
		return
	}

	body, _ := io.ReadAll(res.Body)
	if body == nil {
		ctx.JSON(res.StatusCode, body)
		return
	}

	ctx.JSON(res.StatusCode, body)
}

// GetDocument implements IDatabase.
func (service *EsService) GetDocument() {
	var authorReq *models.Author
	var res *esapi.Response
	var err error
	var success bool

	ctx := service.Ctx
	es := service.Es

	err = ctx.ShouldBindJSON(&authorReq)
	if err != nil || authorReq == nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	if !isTypedApiFlowEnabled {
		res, err = es.Get(authorReq.Details.Index, authorReq.Details.DocumentID)
	} else {
		req := esapi.GetRequest{
			Index:      authorReq.Details.Index,
			DocumentID: authorReq.Details.DocumentID,
		}

		res, err = req.Do(ctx, es)
	}

	success, err = ValidateResponse(res, err)
	if !success || err != nil {
		ctx.JSON(res.StatusCode, err)
		return
	}

	body, _ := io.ReadAll(res.Body)
	if body == nil {
		ctx.JSON(res.StatusCode, nil)
		return
	}

	ctx.JSON(res.StatusCode, body)
}

// UpdateDocument implements IDatabase.
func (service *EsService) UpdateDocument() {
	var authorReq *models.Author
	var res *esapi.Response
	var err error
	var success bool

	ctx := service.Ctx
	es := service.Es

	err = ctx.ShouldBindJSON(&authorReq)
	if err != nil || authorReq == nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	exists, err := DocumentExists(ctx, es, authorReq)
	if err != nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	if !exists {
		fmt.Printf("Document with ID %s doesn't exist. Skipping updation.\n", authorReq.Details.DocumentID)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	updateDoc := map[string]interface{}{
		"doc": map[string]interface{}{
			"name":       authorReq.Name,
			"page_count": authorReq.PageCount,
		},
	}

	jsonUpdateDoc, err := json.Marshal(updateDoc)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	if !isTypedApiFlowEnabled {
		req := esapi.UpdateRequest{
			Index:      authorReq.Details.Index,
			DocumentID: authorReq.Details.DocumentID,
			Body:       bytes.NewReader(jsonUpdateDoc),
		}

		res, err = req.Do(context.Background(), es)
	} else {
		res, err = es.Update(authorReq.Details.Index, authorReq.Details.DocumentID, bytes.NewReader(jsonUpdateDoc))

	}

	defer res.Body.Close()
	success, err = ValidateResponse(res, err)
	if !success || err != nil {
		ctx.JSON(res.StatusCode, err)
		return
	}

	body, _ := io.ReadAll(res.Body)
	if body == nil {
		ctx.JSON(res.StatusCode, nil)
		return
	}

	ctx.JSON(res.StatusCode, body)
}

// Check if document exists
func DocumentExists(ctx *gin.Context, es *elasticsearch.Client, authorReq *models.Author) (bool, error) {

	var resp *esapi.Response
	var err error
	var success bool

	if !isTypedApiFlowEnabled {
		resp, err = es.Exists("search-test", authorReq.Details.DocumentID)
	} else {
		req := esapi.ExistsRequest{
			Index:      authorReq.Details.Index,
			DocumentID: authorReq.Details.DocumentID,
		}

		resp, err = req.Do(context.Background(), es)
	}
	defer resp.Body.Close()

	success, err = ValidateResponse(resp, err)

	if !success || err != nil {
		return success, err
	}

	return success, nil
}

// DeleteDocument implements IDatabase.
func (service *EsService) DeleteDocument() {
	var authorReq *models.Author
	var res *esapi.Response
	var err error
	var success bool

	ctx := service.Ctx
	es := service.Es

	err = ctx.ShouldBindJSON(&authorReq)
	if err != nil || authorReq == nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	exists, err := DocumentExists(ctx, es, authorReq)
	if err != nil {
		log.Fatalf("Error checking if document exists: %s", err)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	if !exists {
		fmt.Printf("Document with ID %s doesn't exist. Skipping updation.\n", authorReq.Details.DocumentID)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	if !isTypedApiFlowEnabled {
		res, err = es.Delete(authorReq.Details.Index, authorReq.Details.DocumentID)
	} else {
		req := esapi.DeleteRequest{
			Index:      authorReq.Details.Index,
			DocumentID: authorReq.Details.DocumentID,
		}

		res, err = req.Do(context.Background(), es)
	}

	defer res.Body.Close()
	success, err = ValidateResponse(res, err)
	if !success || err != nil {
		ctx.JSON(res.StatusCode, err)
		return
	}

	body, _ := io.ReadAll(res.Body)
	if body == nil {
		ctx.JSON(res.StatusCode, nil)
		return
	}

	ctx.JSON(res.StatusCode, body)
}

func ValidateResponse(resp *esapi.Response, err error) (bool, error) {
	if resp == nil || (resp.StatusCode >= 200 && resp.StatusCode <= 300) || err != nil {
		return false, err
	}

	return true, nil
}
