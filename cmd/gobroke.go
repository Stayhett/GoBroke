package main

import (
	"GoBroke/pkg/broker"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func dataHandler(context *broker.Configuration, data []byte) {
	var pipeline broker.PipelineProcessor

	switch context.Input.Type {
	case "csv":
		pipeline = &broker.CSVProcessor{
			Output: context.Output,
			Data:   data,
		}
	default:
		fmt.Println("Not CSV")
		fmt.Println(string(data))
		return
	}
	pipeline.Do()
}

func main() {
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
				dataHandler(config, data)
				wg.Done()
			}(&config, &l)
		}
	}
	wg.Wait()
}
