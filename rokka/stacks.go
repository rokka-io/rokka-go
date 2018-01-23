package rokka

import (
	"net/http"
	"time"
)

// ListStacksResponse contains a list of stacks each containing a list of operations.
type ListStacksResponse struct {
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

// ListStacks returns the stacks for the specified organization.
//
// See: https://rokka.io/documentation/references/stacks.html#retrieve-a-stack
func (c *Client) ListStacks(org string) (ListStacksResponse, error) {
	result := ListStacksResponse{}

	req, err := c.NewRequest(http.MethodGet, "/stacks/"+org, nil, nil)
	if err != nil {
		return result, err
	}

	err = c.callJSONResponse(req, &result)
	return result, err
}
