package elasticsearchgolang

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitElasticClient() (*elasticsearch.Client, error) {

	var config *Config
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
		CloudID: config.Credentials.Id,
		APIKey:  config.Credentials.ApiKey,
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

	return es, nil

}
