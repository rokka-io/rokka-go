package rokka

import "net/http"

// StackOptionsResponse contains the available stack options.
type StackOptionsResponse struct {
	Properties map[string]struct {
		Type      string      `json:"type"`
		MinLength *int        `json:"minLength,omitempty"`
		Values    []string    `json:"values,omitempty"`
		Default   interface{} `json:"default,omitempty"`
		Minimum   *int        `json:"minimum,omitempty"`
		Maximum   *int        `json:"maximum,omitempty"`
	} `json:"properties"`
}

// GetStackOptions returns the available stack options definition.
func (c *Client) GetStackOptions() (StackOptionsResponse, error) {
	result := StackOptionsResponse{}

	req, err := c.NewRequest(http.MethodGet, "/stackoptions", nil, nil)
	if err != nil {
		return result, err
	}

	err = c.callJSONResponse(req, &result)
	return result, err
}
