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

const (
	defaultFlushBytes = int(5e+6)
)

func UploadToElasticHandler(index string, data []map[string]interface{}, output Output) {
	esConfig := &elasticsearch.Config{
		Addresses: []string{os.Getenv(output.Host)},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	if output.Key != "" {
		esConfig.APIKey = os.Getenv(output.Key)
	} else {
		esConfig.Username = os.Getenv(output.Username)
		esConfig.Password = os.Getenv(output.Password)
	}

	UploadToElastic(index, data, esConfig)
}

func UploadToElastic(index string, data []map[string]interface{}, esConfig *elasticsearch.Config) {
	// Upload to elastic
	es, err := elasticsearch.NewClient(*esConfig)
	if err != nil {
		log.Printf("error creating the client: %s", err)
		return
	}

	bi, err := configureBulkIndexer(es, index)
	if err != nil {
		log.Printf("error creating the indexer: %s", err)
		return
	}

	for _, doc := range data {
		data, err := json.Marshal(doc)
		if err != nil {
			log.Printf("error marshalling document: %s", err)
			continue
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Index:  index,
				Body:   bytes.NewReader(data),
				Action: "index",
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("error adding document to the indexer: %s", err)
					} else {
						log.Printf("error: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Printf("error adding document to the indexer: %s", err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Printf("Unexpected error: %s", err)
		return
	}

	stats := bi.Stats()
	if stats.NumFailed > 0 {
		log.Printf("Indexed [%d] documents with [%d] errors", stats.NumFlushed, stats.NumFailed)
	} else {
		log.Printf("Successfully indexed [%d] documents", stats.NumFlushed)
	}
}

func configureBulkIndexer(es *elasticsearch.Client, index string) (esutil.BulkIndexer, error) {
	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         index,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    defaultFlushBytes,
		FlushInterval: 30 * time.Second,
	})
}
