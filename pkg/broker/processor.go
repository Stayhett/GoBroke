package broker

import "time"

type Processor struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
	Data   []byte
}

func appendColumn(table *Table, config map[string]interface{}) *Table {
	header := config["header"].(string)
	value := config["value"].(string)

	for i := range *table {
		if i == 0 {
			(*table)[i] = append((*table)[i], header)
		} else {
			(*table)[i] = append((*table)[i], value)
		}
	}
	return table
}

func timestamp(table *Table) *Table {
	return appendColumn(table, map[string]interface{}{
		"header": "@timestamp",
		"value":  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	})
}
