package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListStacks(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/ListStacks.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListStacks("test-org")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
