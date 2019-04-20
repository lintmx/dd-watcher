package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPGet get page body
func HTTPGet(url string) (string, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(url)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return string(body), err
}
