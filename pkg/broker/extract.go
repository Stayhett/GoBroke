package broker

import (
	"fmt"
	"io"
	"net/http"
)

type Extract struct {
	InputType string
	Data      []byte
}

func (e *Extract) FetchData(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	// Read the response body
	e.Data, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}
	return nil
}
