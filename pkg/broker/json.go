package broker

import (
	"encoding/json"
	"fmt"
)

type JSONProcessor struct {
}

func (J JSONProcessor) Do() {

}

func ReadJSON(jsonData []byte) {
	// Define a variable with the map type to store the decoded JSON data
	var data map[string]interface{}

	// Unmarshal the JSON data into the map
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print the decoded data
	PrintNestedJSON(data, "")
}

func PrintNestedJSON(data map[string]interface{}, prefix string) {
	for key, value := range data {
		switch valueType := value.(type) {
		case map[string]interface{}:
			// Nested JSON object
			fmt.Printf("%s%s:\n", prefix, key)
			PrintNestedJSON(valueType, prefix+"  ")
		case []interface{}:
			// JSON array
			fmt.Printf("%s%s:\n", prefix, key)
			PrintJSONArray(valueType, prefix+"  ")
		default:
			// Primitive data type (string, number, etc.)
			fmt.Printf("%s%s: %v\n", prefix, key, value)
		}
	}
}

func PrintJSONArray(array []interface{}, prefix string) {
	for i, item := range array {
		switch itemType := item.(type) {
		case map[string]interface{}:
			// Nested JSON object in array
			fmt.Printf("%s[%d]:\n", prefix, i)
			PrintNestedJSON(itemType, prefix+"  ")
		default:
			// Primitive data type in array
			fmt.Printf("%s[%d]: %v\n", prefix, i, item)
		}
	}
}
