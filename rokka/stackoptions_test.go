package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStackOptions(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/GetStackOptions.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /stackoptions": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStackOptions()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
