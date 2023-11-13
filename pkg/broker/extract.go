package broker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func fetchData(url string, queue chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	// Download data
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading data from %s: %v\n", url, err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Put the data into the queue
	queue <- data
}
