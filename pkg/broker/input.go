package broker

import "errors"

func InputHandler(inputConfig *Input, location string) ([]byte, error) {
	switch inputConfig.Connector {
	case "http":
		return fetchData(location)
	default:
		return nil, errors.New("error unknown input")
	}
}
