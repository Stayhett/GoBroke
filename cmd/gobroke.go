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

	pipeline := pipelineHandler(&config, data)
	table := pipeline.Do()
	if table == nil {
		log.Println("no data found - no upload")
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

	path := "configuration/"

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal("error opening directory:", err)
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			panic(err)
		}
	}(dir)

	// Read the contents of the directory
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal("error reading directory:", err)
	}

	var configurations []broker.Configuration
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".yml") {
			c, err := broker.ConfigurationConstructorFromFile(path + fileInfo.Name())
			if err != nil {
				log.Println("error constructing configuration:", err)
				continue
			}
			configurations = append(configurations, c)
		}
	}

	var wg sync.WaitGroup
	for _, config := range configurations {
		for _, l := range config.Input.Locations {
			wg.Add(1)
			go processLocation(&wg, config, l)
		}
	}
	wg.Wait()
}
