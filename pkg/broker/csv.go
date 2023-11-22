package broker

import "C"
import (
	"bytes"
	"encoding/csv"
	"log"
)

type CSVConfiguration struct {
	Header []string
}

type CSVSchema struct {
	Mapping   map[string]string `yaml:"mapping"`
	Delimiter string            `yaml:"delimiter"`
	Separator string            `yaml:"separator"`
	Header    string            `yaml:"header"`
	FirstLine bool              `yaml:"FirstLine"`
}

type CSV [][]string

type CSVProcessor struct {
	Output
	CSVConfiguration
	Data   []byte
	Header []string
}

func (p *CSVProcessor) Do() {
	// TODO: Set Header if not done
	data, err := parseCSV(p.Data)
	if err != nil {
		log.Print(err)
		return
	}

	outputData, err := LoadDataHandler(data, p.Output.Type)
	if err != nil {
		log.Print(err)
		return
	}

	err = LoadHandler(outputData, p.Output)
	if err != nil {
		log.Print(err)
		return
	}
}

// parseCSV is a utility function
func parseCSV(data []byte) (CSV, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	return reader.ReadAll()
}
