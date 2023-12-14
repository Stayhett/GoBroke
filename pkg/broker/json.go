package broker

import (
	"encoding/json"
	"fmt"
)

type JSONProcessor struct {
	Output
	Data   []byte
	Schema []string
}

func (J JSONProcessor) Do() {
	readJSON(J.Data)
}

func (J JSONProcessor) GetData() []byte {
	return J.Data
}

func (J JSONProcessor) GetOutput() Output {
	return J.Output
}

func readJSON(jsonData []byte) {
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

	for key, value := range flattenJSON(data, "") {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}
}

// flattenJSON flattens JSON File but only with strings
func flattenJSON(data map[string]interface{}, parentKey string) map[string]interface{} {
	flat := make(map[string]interface{})

	for key, value := range data {
		newKey := fmt.Sprintf("%s%s", parentKey, key)

		switch valueType := value.(type) {
		case map[string]interface{}:
			// Recursively flatten nested JSON
			nestedFlat := flattenJSON(valueType, newKey+".")
			for nestedKey, nestedValue := range nestedFlat {
				flat[nestedKey] = nestedValue
			}
		default:
			// For non-nested values, add them to the flat map
			flat[newKey] = value
		}
	}

	return flat
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
