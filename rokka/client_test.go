package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestValidAPIKey(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetValidAPIKey.json", http.StatusOK)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL, APIKey: "test"})

	ok, err := c.ValidAPIKey()
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("Expected to have a valid API key")
	}
}

func TestInvalidAPIKey(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/GetInvalidAPIKey.json", http.StatusForbidden)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL, APIKey: "test"})

	ok, err := c.ValidAPIKey()
	if err != nil {
		t.Error(err)
	}

	if ok {
		t.Error("Expected to not have a valid API key")
	}
}
