package main

import (
	"GoBroke/pkg/broker"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

//func loadHandler

func dataHandler(context *broker.Configuration, data []byte) {
	switch context.InputType {
	case "csv":
		fmt.Println(string(data))
	default:
		fmt.Println("Not CSV")
		fmt.Println(string(data))
	}
	//loadHandler(res)
}

func main() {
	path := "configuration/"
	var wg sync.WaitGroup

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	// Read the contents of the directory
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	var configurations []broker.Configuration
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".yml") {
			err, c := broker.ConstructConfigurationFromFile(path + fileInfo.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			configurations = append(configurations, c)
		}
	}

	// Working
	for _, config := range configurations {
		// Do pipelining for every location
		for _, l := range config.Input {
			wg.Add(1)
			go func(config *broker.Configuration, lPtr *string) {
				data, err := broker.FetchData(*lPtr)
				if err != nil {
					log.Println(err)
				}
				dataHandler(config, data)
			}(&config, &l)
		}
	}
	wg.Wait()
}
