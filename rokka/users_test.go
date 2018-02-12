package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestCreateUser(t *testing.T) {
	ts := test.NewMockAPI(test.Routes{"POST /users": test.Response{http.StatusOK, "./fixtures/CreateUser.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.CreateUser("test", "test@example.org")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
