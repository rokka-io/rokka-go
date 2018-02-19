package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStats(t *testing.T) {
	org := "test"
	r := test.NewResponse(http.StatusOK, "./fixtures/GetStats.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /stats/" + org: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStats(org, GetStatsOptions{})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
