package broker

import (
	"errors"
	"io"
	"net/http"
	"os"
)

type Connector interface {
	connect() ([]byte, error)
}

func InputHandler(inputConfig *Input, location string) ([]byte, error) {
	var connector Connector
	var err error
	switch inputConfig.Connector {
	case "http":
		connector = NewHTTPConnector(&ConfigHTTPConnector{
			location,
		})
	case "shadowserver":
		connector = NewShadowServerConnector(&ConfigShadowServerConnector{
			secretKey: []byte(inputConfig.IntegrityKey),
			apikey:    inputConfig.Key,
			location:  location,
		})
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

func (I *Input) GetEnvs() {
	I.Key = os.Getenv(I.Key)
	I.IntegrityKey = os.Getenv(I.IntegrityKey)
}
