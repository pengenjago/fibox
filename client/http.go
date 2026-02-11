package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v3/log"
)

// HTTPClientConfig is configuration for HTTP client
type HTTPClientConfig struct {
	BaseURL          string
	Timeout          time.Duration
	Headers          map[string]string
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	Debug            bool
}

// HTTPClient is a wrapper for resty client
type HTTPClient struct {
	client *resty.Client
}

// NewHTTPClient creates a new HTTP client with the given configuration
func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	client := resty.New()

	// Set base URL if provided
	if config.BaseURL != "" {
		client = client.SetBaseURL(config.BaseURL)
	}

	// Set timeout if provided, otherwise use default 30 seconds
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	client = client.SetTimeout(timeout)

	// Set headers if provided
	if config.Headers != nil {
		client = client.SetHeaders(config.Headers)
	}

	// Set retry count if provided
	if config.RetryCount > 0 {
		client = client.SetRetryCount(config.RetryCount)

		// Set retry wait time if provided, otherwise use default 1 second
		retryWaitTime := config.RetryWaitTime
		if retryWaitTime == 0 {
			retryWaitTime = 1 * time.Second
		}
		client = client.SetRetryWaitTime(retryWaitTime)

		// Set retry max wait time if provided, otherwise use default 30 seconds
		retryMaxWaitTime := config.RetryMaxWaitTime
		if retryMaxWaitTime == 0 {
			retryMaxWaitTime = 30 * time.Second
		}
		client = client.SetRetryMaxWaitTime(retryMaxWaitTime)
	}

	// Enable debug mode if requested
	if config.Debug {
		client = client.SetDebug(true)
	}

	// Set default JSON content type
	client = client.SetHeader("Content-Type", "application/json")

	return &HTTPClient{
		client: client,
	}
}

// Get performs a GET request
func (c *HTTPClient) Get(path string, queryParams map[string]string, result interface{}) error {
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(result).
		Get(path)

	if err != nil {
		log.Errorf("HTTP GET request failed: %v", err)
		return fmt.Errorf("HTTP GET request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP GET request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
		return fmt.Errorf("HTTP GET request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return nil
}

// Post performs a POST request
func (c *HTTPClient) Post(path string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetBody(body).
		SetResult(result).
		Post(path)

	if err != nil {
		log.Errorf("HTTP POST request failed: %v", err)
		return fmt.Errorf("HTTP POST request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP POST request %s returned error status: %d, body: %s", path, resp.StatusCode(), resp.Body())
		return fmt.Errorf("HTTP POST request %s returned error status: %d, body: %s", path, resp.StatusCode(), resp.Body())
	}

	return nil
}

// Put performs a PUT request
func (c *HTTPClient) Put(path string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetBody(body).
		SetResult(result).
		Put(path)

	if err != nil {
		log.Errorf("HTTP PUT request failed: %v", err)
		return fmt.Errorf("HTTP PUT request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP PUT request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
		return fmt.Errorf("HTTP PUT request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return nil
}

// Delete performs a DELETE request
func (c *HTTPClient) Delete(path string, queryParams map[string]string, result interface{}) error {
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(result).
		Delete(path)

	if err != nil {
		log.Errorf("HTTP DELETE request failed: %v", err)
		return fmt.Errorf("HTTP DELETE request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP DELETE request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
		return fmt.Errorf("HTTP DELETE request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return nil
}

// PostForm performs a POST request with form data
func (c *HTTPClient) PostForm(path string, formData map[string]string, result interface{}) error {
	resp, err := c.client.R().
		SetFormData(formData).
		SetResult(result).
		Post(path)

	if err != nil {
		log.Errorf("HTTP POST form request failed: %v", err)
		return fmt.Errorf("HTTP POST form request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP POST form request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
		return fmt.Errorf("HTTP POST form request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return nil
}

// GetRaw performs a GET request and returns the raw response
func (c *HTTPClient) GetRaw(path string, queryParams map[string]string) ([]byte, error) {
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		Get(path)

	if err != nil {
		log.Errorf("HTTP GET raw request failed: %v", err)
		return nil, fmt.Errorf("HTTP GET raw request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP GET raw request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
		return nil, fmt.Errorf("HTTP GET raw request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return resp.Body(), nil
}

// PostRaw performs a POST request and returns the raw response
func (c *HTTPClient) PostRaw(path string, body interface{}) ([]byte, error) {
	resp, err := c.client.R().
		SetBody(body).
		Post(path)

	if err != nil {
		log.Errorf("HTTP POST raw request failed: %v", err)
		return nil, fmt.Errorf("HTTP POST raw request failed: %w", err)
	}

	if resp.IsError() {
		log.Errorf("HTTP POST raw request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
		return nil, fmt.Errorf("HTTP POST raw request returned error status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return resp.Body(), nil
}

// SetHeader sets a header for the client
func (c *HTTPClient) SetHeader(key, value string) {
	c.client.SetHeader(key, value)
}

// SetAuthToken sets the authorization token for the client
func (c *HTTPClient) SetAuthToken(token string) {
	c.client.SetAuthToken(token)
}

// SetBasicAuth sets basic authentication for the client
func (c *HTTPClient) SetBasicAuth(username, password string) {
	c.client.SetBasicAuth(username, password)
}

// SetBearerToken sets the bearer token for the client
func (c *HTTPClient) SetBearerToken(token string) {
	c.client.SetAuthToken(token)
}

func (c *HTTPClient) SetDebug(isDebug bool) {
	c.client.SetDebug(isDebug)
}

// GetDefaultHTTPClient returns a default HTTP client with common settings
func GetDefaultHTTPClient(baseURL string) *HTTPClient {
	return NewHTTPClient(HTTPClientConfig{
		BaseURL:          baseURL,
		Timeout:          30 * time.Second,
		RetryCount:       3,
		RetryWaitTime:    1 * time.Second,
		RetryMaxWaitTime: 30 * time.Second,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	})
}

// APIResponse represents a standard API response structure
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ParseAPIResponse parses the API response into a standard structure
func ParseAPIResponse(body []byte) (*APIResponse, error) {
	var response APIResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	return &response, nil
}
