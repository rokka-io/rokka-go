package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestGetStackoptions(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetStackoptions.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.GetStackoptions()
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
