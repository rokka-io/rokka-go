package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStackOptions(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetStackOptions.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStackOptions()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
