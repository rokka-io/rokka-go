package rokka

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type membershipRole string

// RoleRead, RoleWrite and RoleAdmin are membership roles.
const (
	RoleRead  membershipRole = "read"
	RoleWrite membershipRole = "write"
	RoleAdmin membershipRole = "admin"
)

// MembershipRole contains a map of strings to roles.
var MembershipRole = map[string]membershipRole{
	"read":  RoleRead,
	"write": RoleWrite,
	"admin": RoleAdmin,
}

type createMembershipRequest struct {
	Role membershipRole `json:"role"`
}

// CreateMembership creates a membership for a given email address.
//
// See: https://rokka.io/documentation/references/memberships.html#add-user-to-organization
func (c *Client) CreateMembership(org, email string, role membershipRole) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(createMembershipRequest{
		Role: role,
	})
	if err != nil {
		return err
	}
	req, err := c.NewRequest(http.MethodPut, "/organizations/"+org+"/memberships/"+email, b, nil)
	if err != nil {
		return err
	}

	return c.Call(req, nil)
}
