package broker

import (
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
		UploadToElasticHandler("go-broke", outputData, output)

	default:
		return fmt.Errorf("%s is not a known connector", output.Connector)
	}
	return nil
}
