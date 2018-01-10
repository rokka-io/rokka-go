package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStats(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetStats.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStats("test", GetStatsOptions{})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
