package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetOperations(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetOperations.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStackOptions()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
