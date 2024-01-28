package elasticsearchgolang

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitElasticClient() (*elasticsearch.Client, error) {

	cfg := elasticsearch.Config{
		CloudID: "",
		APIKey:  "",
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
