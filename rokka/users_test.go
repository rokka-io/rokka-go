package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestCreateUser(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/CreateUser.json")
	ts := test.NewMockAPI(t, test.Routes{"POST /users": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.CreateUser("test", "test@example.org")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateUserWithoutOrg(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/CreateUser.json")
	ts := test.NewMockAPI(t, test.Routes{"POST /users": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.CreateUserWithoutOrg("test@example.org")
	if err != nil {
		t.Error(err)
	}

	expected := "2692a839-cc59-4c4b-9c39-d44f17ee2094"
	if res.ID != expected {
		t.Errorf("Expected APIKey to be '%s', got: '%s'", expected, res.ID)
	}

	t.Log(res)
}

func TestGetUserID(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/GetUser.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /user": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetUserID()
	if err != nil {
		t.Error(err)
	}

	expected := "SomeUserId"
	if res.ID != expected {
		t.Errorf("Expected APIKey to be '%s', got: '%s'", expected, res.ID)
	}

	t.Log(res)
}
