package main

import (
	"GoBroke/pkg/broker"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	path := "configuration/"

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

	fmt.Println(configurations)
}
