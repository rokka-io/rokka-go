package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestCreateMembership(t *testing.T) {
	org := "test"
	email := "test@example.org"
	ts := test.NewMockAPI(test.Routes{"PUT /organizations/" + org + "/memberships/" + email: test.Response{http.StatusCreated, "./fixtures/CreateMembership.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.CreateMembership(org, email, RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Successfully created membership")
}
