package broker

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Configuration struct {
	InputType  string            `yaml:"type"`
	Input      []string          `yaml:"locations"`
	DataSchema map[string]string `yaml:"schema"`
	Output     string            `yaml:"output"`
	OutputHost string            `yaml:"output.host"`
	OutputType string
}

func ConstructConfigurationFromFile(path string) (error, Configuration) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err, Configuration{}
	}

	var config Configuration
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err, Configuration{}
	}
	// TODO: Add Output Type

	return nil, config
}
