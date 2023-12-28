package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

type Input struct {
	Prefetch     string   `yaml:"prefetch"`
	Type         string   `yaml:"type"`
	Locations    []string `yaml:"locations"`
	Connector    string   `yaml:"connector"`
	Key          string   `yaml:"key"`
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	IntegrityKey string   `yaml:"integrityKey"`
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
	Processors []Processor           `yaml:"processors"`
}

type Configurations []Configuration

func ReadConfigurations(path string) (Configurations, error) {
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

	var configs Configurations
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".yml") {
			c, err := ConfigurationConstructorFromFile(path + fileInfo.Name())
			if err != nil {
				log.Println("error constructing configuration:", err)
				continue
			}
			switch v := c.(type) {
			case Configuration:
				configs = append(configs, v)
			case Configurations:
				configs = append(configs, v...)
			}
		}
	}
	return configs, nil
}

func ConfigurationConstructorFromFile(path string) (interface{}, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	// TODO: Validate Data Scheme
	config.Output.GetEnvs()
	config.Input.GetEnvs()

	if config.Input.Prefetch != "" {
		configs, err := config.PreFetchHandler()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error during preFetching: %v", err))
		}
		if configs != nil {
			return configs, nil
		}
	}
	return config, nil
}

func (C *Configuration) PreFetchHandler() (Configurations, error) {
	var data []byte
	var err error

	// Get data from right connector
	switch C.Input.Connector {
	case "http":
		data, err = FetchData(C.Input.Prefetch)
		// TODO: how to get the input locations
	case "shadowserver":
		conn := shadowServerConnector{
			config: ConfigShadowServerConnector{
				secretKey: []byte(C.Input.IntegrityKey),
				apikey:    C.Input.Key,
			},
			auth:      generateAuth,
			listAPI:   "https://transform.shadowserver.org/api2/reports/list",
			reportAPI: "https://dl.shadowserver.org/",
		}
		data, err = conn.preFetch()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error fetching: %v", err))
		}
		// return with fancy logic

	default:
		return nil, errors.New("no prefetching for this connector available")
	}

	// parse unknown json into map
	var result []map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshalling JSON: %v", err))
	}

	// works at least for shadow server
	var configs Configurations
	for _, report := range result {
		c := Configuration{
			Input:      C.Input,
			DataSchema: C.DataSchema,
			Output:     C.Output,
			Processors: C.Processors,
		}
		c.Input.Locations = []string{report["id"].(string)}
		c.Processors = append(c.Processors, Processor{
			Name: "appendColumn",
			Config: map[string]interface{}{
				"header": "type",
				"value":  report["type"],
			},
		})
		c.Processors = append(c.Processors, Processor{
			Name: "appendColumn",
			Config: map[string]interface{}{
				"header": "report_id",
				"value":  report["id"],
			},
		})
		configs = append(configs, c)
	}

	return configs, nil
}
