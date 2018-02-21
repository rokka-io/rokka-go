package rokka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// StackOptions allows to specify certain settings for a stack. Examples are the compression levels depending on the image format.
type StackOptions map[string]interface{}

// Expression allows to override certain behaviour of a stack based on e.g. DPR size of the requesting client.
type Expression struct {
	Expression string                 `json:"expression"`
	Overrides  map[string]interface{} `json:"overrides"`
}

// Stack specifies the data of a stack.
type Stack struct {
	Organization     string       `json:"organization"`
	Name             string       `json:"name"`
	Created          time.Time    `json:"created"`
	StackOptions     StackOptions `json:"stack_options"`
	StackOperations  Operations   `json:"stack_operations"`
	StackExpressions []Expression `json:"stack_expressions"`
}

// ListStacksResponse contains a list of stacks each containing a list of operations.
type ListStacksResponse struct {
	Items []Stack `json:"items"`
}

// CreateStackRequest specifies the stack to create.
type CreateStackRequest struct {
	Operations  Operations   `json:"operations"`
	Options     StackOptions `json:"options,omitempty"`
	Expressions []Expression `json:"expressions,omitempty"`
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

	err = c.CallJSONResponse(req, &result)

	return result, err
}

// CreateStack allows to create a new stack for the organization.
//
// See: https://rokka.io/documentation/references/stacks.html
func (c *Client) CreateStack(org, name string, stack CreateStackRequest, overwrite bool) (Stack, error) {
	qs := url.Values{}
	if overwrite {
		qs.Set("overwrite", "true")
	}

	result := Stack{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(stack)
	if err != nil {
		return result, err
	}

	req, err := c.NewRequest(http.MethodPut, fmt.Sprintf("/stacks/%s/%s", org, name), b, qs)
	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)
	return result, err
}

// DeleteStack allows to delete an existing stack.
//
// See: https://rokka.io/documentation/references/stacks.html
func (c *Client) DeleteStack(org, name string) error {
	req, err := c.NewRequest(http.MethodDelete, fmt.Sprintf("/stacks/%s/%s", org, name), nil, nil)
	if err != nil {
		return err
	}
	return c.Call(req, nil, nil)
}
