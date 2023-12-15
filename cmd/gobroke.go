package main

import (
	"GoBroke/pkg/broker"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"sync"
)

func pipelineHandler(context *broker.Configuration, data []byte) broker.PipelineProcessor {
	var pipeline broker.PipelineProcessor

	switch context.Input.Type {
	case "csv":
		pipeline = &broker.CSVProcessor{
			Output: context.Output,
			Data:   data,
		}
	case "json":
		pipeline = &broker.JSONProcessor{
			Output: context.Output,
			Data:   data,
		}
	default:
		fmt.Println("not a known type")
		fmt.Println(string(data))
		return nil
	}
	return pipeline
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	path := "configuration/"
	var wg sync.WaitGroup

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {

		}
	}(dir)

	// Read the contents of the directory
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	var configurations []broker.Configuration
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".yml") {
			c, err := broker.ConfigurationConstructorFromFile(path + fileInfo.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			configurations = append(configurations, c)
		}
	}

	// Working
	for _, config := range configurations {
		// Do pipeline for every location
		for _, l := range config.Input.Locations {
			wg.Add(1)
			go func(config *broker.Configuration, lPtr *string) {
				data, err := broker.FetchData(*lPtr)
				if err != nil {
					log.Println(err)
				}

				pipeline := pipelineHandler(config, data)
				table := pipeline.Do()

				err = broker.LoadHandler(table, config.Output)
				if err != nil {
					log.Println(err)
				}
				wg.Done()
			}(&config, &l)
		}
	}
	wg.Wait()
}
