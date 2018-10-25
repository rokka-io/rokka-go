package rokka

import (
	"github.com/rokka-io/rokka-go/test"
	"net/http"
	"testing"
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

func TestCreateNewMembershipWithCurrentUser(t *testing.T) {
	org := "test"
	r := test.NewResponse(http.StatusCreated, "./fixtures/CreateMembershipWithCurrentUser.json")
	ts := test.NewMockAPI(t, test.Routes{"POST /organizations/" + org + "/memberships": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	var roles []MembershipRole
	roles = append(roles, RoleAdmin)
	roles = append(roles, RoleWrite)

	res, err := c.CreateNewMembershipWithCurrentUser(org, roles)
	if err != nil {
		t.Fatal(err)
	}

	expected := "jPALkp8ufA8gfgNStnGuJUG2X0oIWBCD"
	if res.APIKey != expected {
		t.Errorf("Expected APIKey to be '%s', got: '%s'", expected, res.APIKey)
	}

	expected = "newUserId"
	if res.UserID != expected {
		t.Errorf("Expected APIKey to be '%s', got: '%s'", expected, res.UserID)
	}

	expectedRole := "write"
	if res.Roles[1] != expectedRole {
		t.Errorf("Expected second Role to be '%s', got: '%s'", expectedRole, res.Roles[1])
	}
	t.Log(res)
}
