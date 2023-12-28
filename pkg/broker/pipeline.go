package broker

import "fmt"

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
	for _, row := range data { // Skip the header row
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
		switch c.name {
		case "appendColumn":
			AppendColumn(t, c.config)
		default:
			fmt.Printf("unknown processor - %s - continue", c.name)
		}
	}
	return nil
}
