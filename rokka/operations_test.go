package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOperations(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/GetOperations.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /operations": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetOperations()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
