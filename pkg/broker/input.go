package broker

import (
	"errors"
	"io"
	"net/http"
)

func InputHandler(inputConfig *Input, location string) ([]byte, error) {
	switch inputConfig.Connector {
	case "http":
		return fetchData(location)
	default:
		return nil, errors.New("error unknown input")
	}
}

func fetchData(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
