package broker

import "C"
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
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
		err := writeBytesToFile(p.Data, "error.csv")
		if err != nil {
			return nil
		}

		return nil
	}
	return data
}

// parseCSV is a utility function
func parseCSV(data []byte) (Table, error) {
	//ZERO WIDTH NO-BREAK SPACE -> Go reader con not work with it
	hasBOM := len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF

	var reader *csv.Reader
	if hasBOM {
		reader = csv.NewReader(strings.NewReader(string(data[3:])))
	} else {
		reader = csv.NewReader(strings.NewReader(string(data)))
	}

	return reader.ReadAll()
}

func writeBytesToFile(data []byte, filename string) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
