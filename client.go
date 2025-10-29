package main

import (
	"fmt"
	"io"
	"net/http"
)

// fetches the content of a URL, reads it into a string and returns it
func FetchURLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response code fetching URL. status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}
