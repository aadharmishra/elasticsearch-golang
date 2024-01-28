package elasticsearchgolang

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

// Check if document exists
func DocumentExists(ctx *gin.Context, es *elasticsearch.Client, authorReq *Author) (bool, error) {

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
