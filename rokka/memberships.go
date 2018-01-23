package rokka

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MembershipRole string

// RoleRead, RoleWrite and RoleAdmin are membership roles.
const (
	RoleRead  MembershipRole = "read"
	RoleWrite MembershipRole = "write"
	RoleAdmin MembershipRole = "admin"
)

type createMembershipRequest struct {
	Role MembershipRole `json:"role"`
}

// CreateMembership creates a membership for a given email address.
//
// See: https://rokka.io/documentation/references/memberships.html#add-user-to-organization
func (c *Client) CreateMembership(org, email string, role MembershipRole) error {
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

	return c.callJSONResponse(req, nil)
}
