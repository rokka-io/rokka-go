package rokka

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Config Config
}

// HTTPRequester is an interface defining the Do function.
// http.Client is automatically implementing that interface.
type HTTPRequester interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	APIAddress string
	APIVersion string
	APIKey     string
	Verbose    bool
	HTTPClient HTTPRequester
}

type StatusCodeError struct {
	StatusCode int
	Body       []byte
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
		Config: *config,
	}
}

func (c *Client) Call(req *http.Request, v interface{}) error {
	req.Header.Add("Api-Version", c.Config.APIVersion)
	req.Header.Add("Accept", "application/json")

	if len(c.Config.APIKey) != 0 {
		req.Header.Add("Api-Key", c.Config.APIKey)
	}

	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.Config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode < 400 {
		return decoder.Decode(&v)
	}

	errorBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return StatusCodeError{
		resp.StatusCode,
		errorBody,
	}
}

func (c *Client) NewRequest(method, path string, body io.Reader, query map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.Config.APIAddress+path, body)

	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, err
}

func (c *Client) ValidAPIKey() (bool, error) {
	if len(c.Config.APIKey) == 0 {
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
