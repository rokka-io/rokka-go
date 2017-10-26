package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOrganization(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetOrganization.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetOrganization("test")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
