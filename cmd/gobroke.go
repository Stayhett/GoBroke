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

	// Working
	for _, config := range configurations {
		// Do pipelining
		var extracts []broker.Extract

		// TODO: Download Stuff
		for _, l := range config.Input {
			e := broker.Extract{InputType: config.InputType}
			err = e.FetchData(l)
			if err != nil {
				log.Println(err)
			}
			extracts = append(extracts, e)
		}

		// TODO: Parse Stuff
		// TODO: Upload Stuff
		fmt.Println("Names:")
		for _, e := range extracts {
			fmt.Println(e.InputType)
			fmt.Println(string(e.Data))
		}
	}
}
