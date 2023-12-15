package broker

import (
	"fmt"
	"os"
)

func LoadHandler(data Table, output Output) error {
	switch output.Connector {
	case "elasticsearch":
		outputData, err := TableToMaps(data)
		if err != nil {
			return err
		}
		uploadToElasticHandler(outputData, output)
	default:
		return fmt.Errorf("%s is not a known connector", output.Connector)
	}
	return nil
}

func (O *Output) GetEnvs() {
	O.Key = os.Getenv(O.Key)
	O.Username = os.Getenv(O.Username)
	O.Password = os.Getenv(O.Password)
}
