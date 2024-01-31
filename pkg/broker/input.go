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
		connector = NewHTTPConnector(
			&ConfigHTTPConnector{location},
			createNetworkConfig(inputConfig.Proxy))
	case "shadowserver":
		connector = NewShadowServerConnector(&ConfigShadowServerConnector{
			secretKey: []byte(inputConfig.IntegrityKey),
			apikey:    inputConfig.Key,
			location:  location},
			createNetworkConfig(inputConfig.Proxy))
	case "file":
		connector = NewFileConnector(location)
	default:
		return nil, errors.New("unknown input connector")
	}
	if err != nil {
		return nil, err
	}
	return connector.connect()
}

func NewFileConnector(location string) Connector {
	return FileConnector{location}
}

type FileConnector struct {
	path string
}

func (c FileConnector) connect() ([]byte, error) {
	return os.ReadFile(c.path)
}

func NewHTTPConnector(config *ConfigHTTPConnector, networkConfig *http.Transport) Connector {
	return httpConnector{*config, networkConfig}
}

type ConfigHTTPConnector struct {
	url string
}

type httpConnector struct {
	config        ConfigHTTPConnector
	networkConfig *http.Transport
}

func (h httpConnector) connect() ([]byte, error) {
	return FetchData(h.config.url, h.networkConfig)
}

func FetchData(urlString string, networkConfig *http.Transport) ([]byte, error) {
	client := &http.Client{
		Transport: networkConfig,
	}
	req, err := http.NewRequest("GET", urlString, nil)

	response, err := client.Do(req)
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
