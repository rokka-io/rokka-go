package rokka

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client used to communicate with the rokka API.
type Client struct {
	config Config
}

// GetConfig can be used for reading out the configuration of client.
func (c Client) GetConfig() Config {
	return c.config
}

// HTTPRequester is an interface defining the Do function.
// http.Client is automatically implementing that interface.
type HTTPRequester interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config contains configuration for Client.
type Config struct {
	APIAddress string
	APIVersion string
	APIKey     string
	Verbose    bool
	HTTPClient HTTPRequester
}

// StatusCodeError satifies the Error interface and is returned when a response contains a status code >= 400.
type StatusCodeError struct {
	StatusCode int
	Body       []byte
}

// Error creates an error string.
func (e StatusCodeError) Error() string {
	return fmt.Sprintf("rokka: Status Code %d", e.StatusCode)
}

// DefaultConfig is used when calling NewClient with not all config options set.
func DefaultConfig() *Config {
	return &Config{
		APIAddress: "https://api.rokka.io",
		APIVersion: "1",
		APIKey:     "",
		HTTPClient: &http.Client{},
	}
}

// NewClient returns a new client
func NewClient(config *Config) (c *Client) {
	defConfig := DefaultConfig()

	if len(config.APIAddress) == 0 {
		config.APIAddress = defConfig.APIAddress
	}

	if len(config.APIVersion) == 0 {
		config.APIVersion = defConfig.APIVersion
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

// Call executes an HTTP request. It automatically adds necessary headers and decodes the JSON body into `v`.
// If the response contains a status code >= 400 a StatusCodeError is returned.
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return StatusCodeError{
			resp.StatusCode,
			body,
		}
	}
	// handle empty responses
	if len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, &v)
}

// NewRequest constructs a new http.Request used for executing using Call.
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

// ValidAPIKey can be used to check if the API key is valid. It will execute a request to `/` which is an undocumented API.
// This function only returns true if there has been no error and the status code is < 400.
func (c *Client) ValidAPIKey() (bool, error) {
	if len(c.config.APIKey) == 0 {
		return false, errors.New("API key must be set")
	}

	req, err := c.NewRequest(http.MethodGet, "/", nil, nil)
	if err != nil {
		return false, err
	}
	err = c.Call(req, nil)
	if err != nil {
		// only 403 is an expected error code, just return false without the error in this case.
		if err, ok := err.(StatusCodeError); ok && err.StatusCode == http.StatusForbidden {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
