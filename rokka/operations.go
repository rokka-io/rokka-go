package rokka

import "net/http"

//go:generate go run ../cmd/gen/operations.go

// OperationsResponse contains the available stack options.
type OperationsResponse map[string]map[string]interface{}

// GetOperations returns the available stack options definition.
//
// See: https://rokka.io/documentation/references/operations.html
func (c *Client) GetOperations() (OperationsResponse, error) {
	result := OperationsResponse{}

	req, err := c.NewRequest(http.MethodGet, "/operations", nil, nil)
	if err != nil {
		return result, err
	}

	err = c.callJSONResponse(req, &result)
	return result, err
}
