package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStats(t *testing.T) {
	org := "test"
	ts := test.NewMockAPI(test.Routes{"GET /stats/" + org: test.Response{http.StatusOK, "./fixtures/GetStats.json", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStats(org, GetStatsOptions{})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
