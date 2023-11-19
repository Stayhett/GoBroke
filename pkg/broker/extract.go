package broker

import (
	"io"
	"net/http"
)

func FetchData(url string) ([]byte, error) {
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
