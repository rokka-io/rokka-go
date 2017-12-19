package rokka

import "net/http"

// OperationsResponse contains the available stack options.
type OperationsResponse map[string]interface{}

// GetOperations returns the available stack options definition.
//
// See: https://rokka.io/documentation/references/operations.html
func (c *Client) GetOperations() (OperationsResponse, error) {
	result := OperationsResponse{}

	req, err := c.NewRequest(http.MethodGet, "/operations", nil, nil)
	if err != nil {
		return result, err
	}

	err = c.Call(req, &result)
	return result, err
}
