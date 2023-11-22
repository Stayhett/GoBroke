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

func (o *Output) determineOutputType() error {
	o.Type = "json"
	return nil
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

	err = config.Output.determineOutputType()
	if err != nil {
		return Configuration{}, err
	}
	// TODO: Validate Data Scheme
	return config, nil
}
