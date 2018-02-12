package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOrganization(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"GET /organizations/" + org: test.Response{http.StatusOK, "./fixtures/GetOrganization.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetOrganization(org)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateOrganization(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"PUT /organizations/" + org: test.Response{http.StatusOK, "./fixtures/CreateOrganization.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.CreateOrganization(org, "info@example.com", "Dev Environment")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
