package rokka

import "time"

type StacksListResponse struct {
	Items []struct {
		Organization    string                 `json:"organization"`
		Name            string                 `json:"name"`
		Created         time.Time              `json:"created"`
		StackOptions    map[string]interface{} `json:"stack_options"`
		StackOperations []struct {
			Name    string
			Options map[string]interface{}
		} `json:"stack_operations"`
	} `json:"items"`
}

// ListStacks returns the stacks for the specified organizaton.
func (c *Client) ListStacks(org string) (StacksListResponse, error) {
	result := StacksListResponse{}

	req, err := c.NewRequest("GET", "/stacks/"+org, nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
