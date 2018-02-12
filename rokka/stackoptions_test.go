package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStackOptions(t *testing.T) {
	ts := test.NewMockAPI(test.Routes{"GET /stackoptions": test.Response{http.StatusOK, "./fixtures/GetStackOptions.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStackOptions()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
