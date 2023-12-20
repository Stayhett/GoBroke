package broker

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func NewShadowServerConnector(config *ConfigShadowServerConnector) Connector {
	return shadowServerConnector{
		config:    *config,
		auth:      generateAuth,
		listAPI:   "https://transform.shadowserver.org/api2/reports/list",
		reportAPI: "https://dl.shadowserver.org/",
	}
}

func generateAuth(message []byte, secret []byte) ([]byte, error) {
	hasher := hmac.New(sha256.New, secret)

	hasher.Write(message)
	hmacResult := hasher.Sum(nil)
	return hmacResult, nil
}

type ConfigShadowServerConnector struct {
	secretKey []byte
	apikey    string
}

// Shadow server Input Connector
type shadowServerConnector struct {
	config    ConfigShadowServerConnector
	auth      func(message []byte, secret []byte) ([]byte, error)
	listAPI   string
	reportAPI string
	location  string
}

func (s shadowServerConnector) connect() ([]byte, error) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	jsonData, err := json.Marshal(map[string]interface{}{
		"apikey": s.config.apikey,
		"query": map[string]interface{}{
			"data": yesterday.Format("2006-01-02"),
		},
	})
	if err != nil {
		return nil, err
	}

	return s.callAPI(jsonData, s.reportAPI+s.location)
}

func (s shadowServerConnector) preFetch() ([]byte, error) {
	jsonData, err := json.Marshal(map[string]interface{}{"apikey": s.config.apikey})
	if err != nil {
		return nil, err
	}

	return s.callAPI(jsonData, s.listAPI)
}

func (s shadowServerConnector) callAPI(data []byte, url string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	integrity, err := s.auth(data, s.config.secretKey)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Hmac2", hex.EncodeToString(integrity))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}
