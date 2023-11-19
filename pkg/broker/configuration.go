package broker

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Input struct {
	Type      string   `yaml:"type"`
	Locations []string `yaml:"locations"`
	Auth      string
}

type Output struct {
	Connector string `yaml:"connector"`
	Host      string `yaml:"host"`
	Type      string
}

func (o *Output) determineOutputType() {
	// TODO: To Implement
}

type DataSchema struct {
	Mapping   map[string]string `yaml:"mapping"`
	Delimiter string            `yaml:"delimiter"`
	Separator string            `yaml:"separator"`
	Header    []string          `yaml:"header"`
}

type Configuration struct {
	Input      Input      `yaml:"input"`
	Output     Output     `yaml:"output"`
	DataSchema DataSchema `yaml:"schema"`
}

func ConfigurationConstructorFromFile(path string) (error, Configuration) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err, Configuration{}
	}

	var config Configuration
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		return err, Configuration{}
	}

	config.Output.determineOutputType()
	// TODO: Validate Data Scheme
	return nil, config
}
