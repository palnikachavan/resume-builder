package database

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

// ESClient is the global Elasticsearch client
var ESClient *elasticsearch.Client

// InitElasticsearch initializes the Elasticsearch client
func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ES_URL"),
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Check connection
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response from Elasticsearch: %s", err)
	}

	defer res.Body.Close()
	fmt.Println("Elasticsearch initialized:", res)
	ESClient = client
}
