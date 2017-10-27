package rokka

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	config Config
}

type Config struct {
	APIAddress string
	APIVersion string
	APIKey     string
	HTTPClient *http.Client
}

type StatusCodeError struct {
	StatusCode int
	Body       ErrorResponse
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e StatusCodeError) Error() string {
	return fmt.Sprintf("rokka: Status Code %d", e.StatusCode)
}

func DefaultConfig() *Config {
	return &Config{
		APIAddress: "https://api.rokka.io",
		APIVersion: "1",
		APIKey:     "",
		HTTPClient: &http.Client{},
	}
}

func NewClient(config *Config) (c *Client) {
	defConfig := DefaultConfig()

	if len(config.APIAddress) == 0 {
		config.APIAddress = defConfig.APIAddress
	}

	if len(config.APIKey) == 0 {
		config.APIKey = defConfig.APIKey
	}

	if config.HTTPClient == nil {
		config.HTTPClient = defConfig.HTTPClient
	}

	return &Client{
		config: *config,
	}
}

func (c *Client) Call(req *http.Request, v interface{}) error {
	req.Header.Add("Api-Version", c.config.APIVersion)
	req.Header.Add("Accept", "application/json")

	if len(c.config.APIKey) != 0 {
		req.Header.Add("Api-Key", c.config.APIKey)
	}

	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode < 400 {
		return decoder.Decode(&v)
	}

	errorBody := ErrorResponse{}
	err = decoder.Decode(&errorBody)
	if err != nil {
		return err
	}

	return StatusCodeError{
		resp.StatusCode,
		errorBody,
	}
}

func (c *Client) NewRequest(method, path string, body io.Reader, query map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.config.APIAddress+path, body)

	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, err
}
