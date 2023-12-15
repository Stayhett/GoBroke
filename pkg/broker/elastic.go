package broker

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func UploadToElasticHandler(index string, data []map[string]interface{}, output Output) {
	var esConfig elasticsearch.Config
	if output.Key != "" {
		esConfig = elasticsearch.Config{
			Addresses: []string{os.Getenv(output.Host)},
			APIKey:    os.Getenv(output.Key),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		esConfig = elasticsearch.Config{
			Addresses: []string{os.Getenv(output.Host)},
			Username:  os.Getenv(output.Username),
			Password:  os.Getenv(output.Password),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	UploadToElastic(index, data, &esConfig)
}

func UploadToElastic(index string, data []map[string]interface{}, esConfig *elasticsearch.Config) {
	// Upload to elastic
	es, err := elasticsearch.NewClient(*esConfig)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         index,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    int(5e+6),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	for _, doc := range data {
		data, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Error marshalling document: %s", err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Index:  index,
				Body:   bytes.NewReader(data),
				Action: "index",
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Error adding document to the indexer: %s", err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
	stats := bi.Stats()
	if stats.NumFailed > 0 {
		log.Fatalf("Indexed [%d] documents with [%d] errors", stats.NumFlushed, stats.NumFailed)
	} else {
		log.Printf("Successfully indexed [%d] documents", stats.NumFlushed)
	}
}
