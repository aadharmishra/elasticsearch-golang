package elasticsearchgolang

import (
	"io"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// GetDocument implements IDatabase.
func (service *esService) GetDocument() {
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
