package broker

type Processor struct {
	name   string
	config map[string]interface{}
	data   []byte
}

func AppendColumn(table *Table, config map[string]interface{}) *Table {
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
