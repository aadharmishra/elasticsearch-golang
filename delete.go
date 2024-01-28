package elasticsearchgolang

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// DeleteDocument implements IDatabase.
func (service *esService) DeleteDocument() {
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
