package rokka

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// CreateUserResponse contains information about a newly created organization and user
type CreateUserResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

type createUserRequest struct {
	Organization string `json:"organization"`
	Email        string `json:"email"`
}

// CreateUser creates a new organization and a user with the given email address.
//
// See: https://rokka.io/documentation/references/users.html#create-a-user
func (c *Client) CreateUser(org, email string) (CreateUserResponse, error) {
	result := CreateUserResponse{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createUserRequest{
		org,
		email,
	})
	if err != nil {
		return result, err
	}
	req, err := c.NewRequest(http.MethodPost, "/users", b, nil)
	if err != nil {
		return result, err
	}

	err = c.callJSONResponse(req, &result)
	return result, err
}
