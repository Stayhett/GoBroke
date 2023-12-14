package broker

import (
	"encoding/json"
	"fmt"
)

func LoadHandler(data Table, output Output) error {
	switch output.Connector {
	case "elasticsearch":
		fmt.Printf("try upload data to elasticsearch")
		outputData, err := TableToMaps(data)
		if err != nil {
			return err
		}
		UploadToElastic("go-broke", outputData)

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
