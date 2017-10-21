package rokka

type Property struct {
	Type      string      `json:"type"`
	MinLength *int        `json:"minLength,omitempty"`
	Values    []string    `json:"values,omitempty"`
	Default   interface{} `json:"default,omitempty"`
	Minimum   *int        `json:"minimum,omitempty"`
	Maximum   *int        `json:"maximum,omitempty"`
}

type StackoptionsResponse struct {
	Properties map[string]Property `json:"properties"`
}

func (c *Client) GetStackoptions() (StackoptionsResponse, error) {
	result := StackoptionsResponse{}

	req, err := c.NewRequest("GET", "/stackoptions", nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
