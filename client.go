package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// Get performs a GET request to the specified endpoint
func (c *Client) Get(endpoint string) (string, error) {
	url := endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Validate JSON response
	if !isValidJSON(string(body)) {
		return "", fmt.Errorf("malformed response: invalid JSON")
	}

	return string(body), nil
}

// NewClient creates a new API client
func NewClient(baseURL string) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			for key, val := range via[0].Header {
				req.Header[key] = val
			}
			return nil
		},
	}

	SetGlobalUrl(baseURL)

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: client,
	}, nil
}

// isValidJSON checks if a string is valid JSON
func isValidJSON(str string) bool {
	var js interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}
