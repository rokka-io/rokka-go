package rokka

import (
	"net/http"
	"testing"
	"time"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOrganization(t *testing.T) {
	org := "test"
	r := test.NewResponse(http.StatusOK, "./fixtures/GetOrganization.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /organizations/" + org: r})
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
	r := test.NewResponse(http.StatusOK, "./fixtures/CreateOrganization.json")
	ts := test.NewMockAPI(t, test.Routes{"PUT /organizations/" + org: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.CreateOrganization(org, "info@example.com", "Dev Environment")
	if err != nil {
		t.Error(err)
	}
	expected, _ := time.Parse(time.RFC3339, "2018-02-14T07:47:46Z")
	if res.Created.String() != expected.String() {
		t.Errorf("Expected created to be '%s', got: '%s'", expected, res.Created.String())
	}
	if res.Limit != nil {
		t.Errorf("Expected limit to be nil, got: %#v", res.Limit)
	}

	t.Log(res)
}
