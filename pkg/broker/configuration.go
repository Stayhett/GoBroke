package broker

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Input struct {
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
	Location  string `yaml:"index"`
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

	if err != nil {
		return Configuration{}, err
	}
	// TODO: Validate Data Scheme
	return config, nil
}
