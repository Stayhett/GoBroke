package broker

import (
	"errors"
	"io"
	"net/http"
)

type Connector interface {
	connect() ([]byte, error)
}

func InputHandler(inputConfig *Input, location string) ([]byte, error) {
	var connector Connector
	var err error
	switch inputConfig.Connector {
	case "http":
		connector, err = NewHTTPConnector(&ConfigHTTPConnector{
			location,
		}), nil
	case "shadowserver":
		connector, err = NewShadowServerConnector(&ConfigShadowServerConnector{}), nil
	default:
		return nil, errors.New("unknown input connector")
	}
	if err != nil {
		return nil, err
	}
	return connector.connect()
}

func NewHTTPConnector(config *ConfigHTTPConnector) Connector {
	return httpConnector{*config}
}

type ConfigHTTPConnector struct {
	url string
}

type httpConnector struct {
	config ConfigHTTPConnector
}

func (h httpConnector) connect() ([]byte, error) {
	return FetchData(h.config.url)
}

func FetchData(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
