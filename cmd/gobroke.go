package main

import (
	"GoBroke/pkg/broker"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

func pipelineHandler(context *broker.Configuration, data []byte) broker.PipelineProcessor {
	switch context.Input.Type {
	case "csv":
		return &broker.CSVProcessor{
			Output: context.Output,
			Data:   data,
		}
	case "json":
		return &broker.JSONProcessor{
			Output: context.Output,
			Data:   data,
		}
	default:
		fmt.Println("Unknown input type:", context.Input.Type)
		return nil
	}
}

func processLocation(wg *sync.WaitGroup, config broker.Configuration, location string) {
	defer wg.Done()

	data, err := broker.InputHandler(&config.Input, location)
	if err != nil {
		log.Println("error fetching data:", err)
		return
	}

	// Table handling
	pipeline := pipelineHandler(&config, data)
	table := pipeline.Do()
	if table == nil {
		log.Println("no data found - no upload")
		return
	}

	err = table.Process(&config.Processors)
	if err != nil {
		log.Println("error during processing: ", err)
		return
	}

	err = broker.LoadHandler(table, config.Output)
	if err != nil {
		log.Println("error loading data:", err)
	}
	return
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file:", err)
	}

	configurations, err := broker.ReadConfigurations("configuration/")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, config := range configurations {
		for _, l := range config.Input.Locations {
			wg.Add(1)
			configCopy := config
			go processLocation(&wg, configCopy, l)
		}
	}
	wg.Wait()
}
