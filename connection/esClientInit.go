package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"esmodule/models"

	"github.com/elastic/go-elasticsearch/v8"
)

var ConnPool map[string]*elasticsearch.Client

func InitElasticClient() (*elasticsearch.Client, error) {

	var config *models.Config
	var err error
	var data []byte

	configFile := "/Users/aadharmishra/Documents/github/elasticsearch-golang/config.json"
	data, err = os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)

	if err != nil || config == nil {
		log.Fatalf("Error getting creds: %s", err)
		return nil, err
	}

	cfg := elasticsearch.Config{
		CloudID: config.EsCredentials.Id,
		APIKey:  config.EsCredentials.ApiKey,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)

		// API Key should have cluster monitoring rights
		infores, err := es.Info()
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
			return nil, err
		}

		fmt.Println(infores)
		return nil, err
	}

	//global connection pool
	ConnPool = make(map[string]*elasticsearch.Client)
	ConnPool["client"] = es

	return es, nil

}
