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
