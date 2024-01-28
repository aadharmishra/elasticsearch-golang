package elasticsearchgolang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Create document
func (service *esService) CreateDocument() {

	var authorReq *Author
	var res *esapi.Response
	var err error
	var success bool

	ctx := service.ctx
	es := service.es

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
