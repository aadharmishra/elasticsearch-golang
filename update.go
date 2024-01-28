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

// UpdateDocument implements IDatabase.
func (service *esService) UpdateDocument() {
	var authorReq *Author
	var res *esapi.Response
	var err error
	var success bool

	ctx := service.ctx
	es := service.es

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
