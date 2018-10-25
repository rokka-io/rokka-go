package rokka

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// MembershipRole is a simple alias to string for the defined constants.
type MembershipRole string

// RoleRead, RoleWrite, RoleAdmin and RoleUpload are membership roles.
const (
	RoleRead   MembershipRole = "read"
	RoleWrite  MembershipRole = "write"
	RoleAdmin  MembershipRole = "admin"
	RoleUpload MembershipRole = "upload"
)

type createMembershipRequest struct {
	Roles []MembershipRole `json:"roles"`
}

// ListMembershipsResponse is a collection of multiple membership objects in Items.
type ListMembershipsResponse struct {
	Items []Membership `json:"items"`
}

// Membership holds all the info for a membership.
type Membership struct {
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	UserID         string   `json:"user_id"`
	Roles          []string `json:"roles"`
	APIKey         string   `json:"api_key"`
}

// CreateMembership creates a membership for a given UUID of the user.
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#assign-a-user-to-an-organization
func (c *Client) CreateMembership(org, userid string, roles []MembershipRole) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createMembershipRequest{
		Roles: roles,
	})
	if err != nil {
		return err
	}
	req, err := c.NewRequest(http.MethodPut, "/organizations/"+org+"/memberships/"+userid, b, nil)
	if err != nil {
		return err
	}

	return c.Call(req, nil, nil)
}

// CreateNewMembershipWithCurrentUser creates a new user object and automatically assign it to an organisation
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#create-a-new-user-object-and-automatically-assign-it-to-an-organisation
func (c *Client) CreateNewMembershipWithCurrentUser(org string, roles []MembershipRole) (Membership, error) {
	result := Membership{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createMembershipRequest{
		Roles: roles,
	})
	if err != nil {
		return result, err
	}
	req, err := c.NewRequest(http.MethodPost, "/organizations/"+org+"/memberships", b, nil)
	if err != nil {
		return result, err
	}
	err = c.CallJSONResponse(req, &result)

	return result, err
}

// ListMembership lists all membership of an organisation
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#list-memberships
func (c *Client) ListMembership(org string) (ListMembershipsResponse, error) {
	result := ListMembershipsResponse{}

	req, err := c.NewRequest(http.MethodGet, "/organizations/"+org+"/memberships", nil, nil)

	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)

	return result, err
}

// ListMembershipForUUID list the membership of an user and organisation
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#list-memberships
func (c *Client) ListMembershipForUUID(org string, uuid string) (Membership, error) {
	result := Membership{}

	req, err := c.NewRequest(http.MethodGet, "/organizations/"+org+"/memberships/"+uuid, nil, nil)

	if err != nil {
		return result, err
	}

	err = c.CallJSONResponse(req, &result)

	return result, err
}

// DeleteMembership deletes a user from an organization
//
// See: https://rokka.io/documentation/references/users-and-memberships.html#remove-a-user-from-an-organization
func (c *Client) DeleteMembership(org string, uuid string) error {
	req, err := c.NewRequest(http.MethodDelete, "/organizations/"+org+"/memberships/"+uuid, nil, nil)

	if err != nil {
		return err
	}

	return c.Call(req, nil, nil)
}
