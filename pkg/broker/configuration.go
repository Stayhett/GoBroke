package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Input struct {
	Prefetch  string   `yaml:"prefetch"`
	Type      string   `yaml:"type"`
	Locations []string `yaml:"locations"`
	Connector string   `yaml:"connector"`
	Key       string   `yaml:"key"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
}

type Output struct {
	Connector string `yaml:"connector"`
	Host      string `yaml:"host"`
	Store     string `yaml:"store"`
	Key       string `yaml:"key"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type PipelineConfiguration struct {
	CSVSchema
}

type Configuration struct {
	Input      Input                 `yaml:"input"`
	Output     Output                `yaml:"output"`
	DataSchema PipelineConfiguration `yaml:"schema"`
}

func ConfigurationConstructorFromFile(path string) (Configuration, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Configuration{}, err
	}

	// TODO: Validate Data Scheme
	config.Output.GetEnvs()
	if config.Input.Prefetch != "" {
		err = config.PreFetchHandler()
		if err != nil {
			return Configuration{}, err
		}
	}
	return config, nil
}

func (C Configuration) PreFetchHandler() error {
	var data []byte
	var err error

	// Get data from right connector
	switch C.Input.Connector {
	case "http":
		data, err = FetchData(C.Input.Prefetch)
		// TODO: how to get the input locations
	case "shadowserver":
		conn := shadowServerConnector{}
		data, err = conn.preFetch()
	default:
		return errors.New("no prefetching for this connector available")
	}

	// parse unknown json into map
	var result []map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return errors.New(fmt.Sprintf("error unmarshalling JSON: %v", err))
	}

	// works at least for shadowserver
	for _, report := range result {
		if value, ok := report["id"]; ok {
			C.Input.Locations = append(C.Input.Locations, value.(string))
		}
	}

	return nil
}
