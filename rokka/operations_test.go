package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOperations(t *testing.T) {
	ts := test.NewMockAPI(test.Routes{"GET /operations": test.Response{http.StatusOK, "./fixtures/GetOperations.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetOperations()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
