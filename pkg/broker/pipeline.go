package broker

import (
	"errors"
	"fmt"
)

type PipelineProcessor interface {
	Do() Table
}

// TODO: Explicit type for json -> json Processor

type Pipeline struct {
	Output
}

type Table [][]string

func TableToMaps(data Table) ([]map[string]interface{}, error) {
	var csvMaps []map[string]interface{}
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	for _, row := range data[1:] {
		csvMap := make(map[string]interface{})
		for i, value := range row {
			header := data[0][i]
			csvMap[header] = value
		}
		csvMaps = append(csvMaps, csvMap)
	}

	return csvMaps, nil
}

func (t *Table) Process(config *[]Processor) error {
	for _, c := range *config {
		switch c.Name {
		case "appendColumn":
			appendColumn(t, c.Config)
		case "timestamp":
			timestamp(t)
		default:
			fmt.Printf("unknown processor - %s - continue\n", c.Name)
		}
	}
	return nil
}
