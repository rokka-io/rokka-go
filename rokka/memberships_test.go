package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestCreateMembership(t *testing.T) {
	org := "test"
	email := "test@example.org"
	r := test.NewResponse(http.StatusCreated, "./fixtures/CreateMembership.json")
	ts := test.NewMockAPI(t, test.Routes{"PUT /organizations/" + org + "/memberships/" + email: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	var roles []MembershipRole
	roles = append(roles, RoleAdmin)
	err := c.CreateMembership(org, email, roles)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Successfully created membership")
}
