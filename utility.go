package elasticsearchgolang

import "github.com/elastic/go-elasticsearch/v8/esapi"

func ValidateResponse(resp *esapi.Response, err error) (bool, error) {
	if resp == nil || resp.StatusCode != 200 || err != nil {
		return false, err
	}

	return true, nil
}
