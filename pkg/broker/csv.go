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

type CSVProcessor struct {
	Output
	CSVConfiguration
	Data   []byte
	Header []string
}

func (p *CSVProcessor) Do() Table {
	// TODO: Set Header if not done
	data, err := parseCSV(p.Data)
	if err != nil {
		log.Printf("error in parse csv: %s", err)
		return nil
	}
	return data
}

// parseCSV is a utility function
func parseCSV(data []byte) (Table, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	return reader.ReadAll()
}
