package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is an HTTP client for the not-env API
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

// Request performs an HTTP request
func (c *Client) Request(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(path string) (*http.Response, error) {
	return c.Request("GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.Request("POST", path, body)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.Request("PUT", path, body)
}

// Patch performs a PATCH request
func (c *Client) Patch(path string, body interface{}) (*http.Response, error) {
	return c.Request("PATCH", path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Request("DELETE", path, nil)
}

// ParseResponse parses a JSON response
func ParseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return fmt.Errorf("%s: %s", errResp.Error, errResp.Message)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

