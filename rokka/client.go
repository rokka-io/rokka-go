package rokka

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var errorAPIKeyMissing = errors.New("API key must be set")

// AnnotatedUnmarshalTypeError is a wrapper for json.UnmarshalTypeError adding the offending JSON body around the offset.
type AnnotatedUnmarshalTypeError struct {
	*json.UnmarshalTypeError
	Content string
}

// Error returns the same error as UnmarshalTypeError.
func (a *AnnotatedUnmarshalTypeError) Error() string {
	return a.UnmarshalTypeError.Error()
}

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
	ImageHost  string
	Verbose    bool
	HTTPClient HTTPRequester
}

// APIError is returned by the API in case of errors.
type APIError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// StatusCodeError satifies the Error interface and is returned when a response contains a status code >= 400.
type StatusCodeError struct {
	Code     int
	APIError *APIError
	Body     []byte
}

// Error creates an error string.
func (e StatusCodeError) Error() string {
	s := fmt.Sprintf("rokka: Status Code %d", e.Code)
	if e.APIError != nil {
		s += fmt.Sprintf(" (%s)", e.APIError.Error.Message)
	}
	return s
}

// DefaultConfig is used when calling NewClient with not all config options set.
func DefaultConfig() *Config {
	return &Config{
		APIAddress: "https://api.rokka.io",
		APIVersion: "1",
		APIKey:     "",
		ImageHost:  "https://{{organization}}.rokka.io",
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

	if len(config.ImageHost) == 0 {
		config.ImageHost = defConfig.ImageHost
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
		rErr := APIError{}
		sErr := StatusCodeError{
			Code: resp.StatusCode,
			Body: body,
		}
		if err := json.Unmarshal(body, &rErr); err != nil {
			return sErr
		}
		sErr.APIError = &rErr
		return sErr
	}
	if err := json.Unmarshal(body, &v); err != nil {
		switch err := err.(type) {
		case *json.UnmarshalTypeError:
			return &AnnotatedUnmarshalTypeError{
				UnmarshalTypeError: err,
				Content:            fmt.Sprintf("%s\n<-->\n%s", body[err.Offset-100:err.Offset], body[err.Offset:err.Offset+100]),
			}
		default:
			return err
		}
	}
	return nil
}

// NewRequest constructs a new http.Request used for executing using Call.
func (c *Client) NewRequest(method, path string, body io.Reader, query url.Values) (*http.Request, error) {
	req, err := http.NewRequest(method, c.config.APIAddress+path, body)

	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	return req, err
}

// ValidAPIKey can be used to check if the API key is valid. It will execute a request to `/` which is an undocumented API.
// This function only returns true if there has been no error and the status code is < 400.
func (c *Client) ValidAPIKey() (bool, error) {
	if len(c.config.APIKey) == 0 {
		return false, errorAPIKeyMissing
	}

	req, err := c.NewRequest(http.MethodGet, "/", nil, nil)
	if err != nil {
		return false, err
	}
	err = c.Call(req, nil)
	if err != nil {
		// only 403 is an expected error code, just return false without the error in this case.
		if err, ok := err.(StatusCodeError); ok && err.Code == http.StatusForbidden {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
