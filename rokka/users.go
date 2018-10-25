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

// CreateUserIDResponse contains information about the user id of a user
type CreateUserIDResponse struct {
	ID string `json:"user_id"`
}

type createUserAndOrgRequest struct {
	Organization string `json:"organization"`
	Email        string `json:"email"`
}

type createUser struct {
	Email string `json:"email"`
}

// CreateUser creates a new organization and a user with the given email address.
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#create-a-user
func (c *Client) CreateUser(org, email string) (CreateUserResponse, error) {
	result := CreateUserResponse{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createUserAndOrgRequest{
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

	err = c.CallJSONResponse(req, &result)
	return result, err
}

// GetUserID returns the user uuid of the logged in user
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#get-the-current-user_id
func (c *Client) GetUserID() (CreateUserIDResponse, error) {
	result := CreateUserIDResponse{}

	req, err := c.NewRequest(http.MethodGet, "/user", nil, nil)
	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)
	return result, err
}

// CreateUserWithoutOrg creates a user with the given email address.
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#create-a-user
func (c *Client) CreateUserWithoutOrg(email string) (CreateUserResponse, error) {
	result := CreateUserResponse{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createUser{
		email,
	})
	if err != nil {
		return result, err
	}
	req, err := c.NewRequest(http.MethodPost, "/users", b, nil)
	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)
	return result, err
}
