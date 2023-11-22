package broker

import (
	"encoding/json"
	"fmt"
)

func LoadDataHandler(data [][]string, outputType string) ([]byte, error) {
	var outputData []byte
	var err error
	switch outputType {
	case "json":
		outputData, err = parseCSVToMaps(data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s is not a known output type", outputType)
	}
	return outputData, nil
}

func LoadHandler(data []byte, output Output) error {
	switch output.Connector {
	case "elasticsearch":
		fmt.Printf("Would upload data to elasticsearch")
		fmt.Println(string(data))
	default:
		return fmt.Errorf("%s is not a known connector", output.Connector)
	}
	return nil
}

func parseCSVToMaps(data [][]string) ([]byte, error) {
	var csvMaps []map[string]interface{}
	for _, row := range data { // Skip the header row
		csvMap := make(map[string]interface{})
		for i, value := range row {
			header := data[0][i]
			csvMap[header] = value
		}
		csvMaps = append(csvMaps, csvMap)
	}
	jsonData, err := json.Marshal(csvMaps)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return nil, err
	}
	return jsonData, nil
}
