package broker

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
