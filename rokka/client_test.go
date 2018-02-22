package rokka

import (
	"net/http"
	"strings"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestCall_JSONError(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/Call_JSONError.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /": r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL, APIKey: "test"})

	req, err := c.NewRequest(http.MethodGet, "/", nil, nil)
	if err != nil {
		panic(err)
	}
	out := struct {
		Width int `json:"width"`
	}{}

	err = c.CallJSONResponse(req, &out)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	rErr, ok := err.(*AnnotatedUnmarshalTypeError)
	if !ok {
		t.Fatalf("Expected '%T', got '%T'", new(AnnotatedUnmarshalTypeError), err)
	}
	expectedErrStrPrefix := "json: cannot unmarshal string"
	if !strings.HasPrefix(rErr.Error(), expectedErrStrPrefix) {
		t.Errorf("Expected error to have prefix '%s', got '%s'", expectedErrStrPrefix, rErr.Error())
	}
	expectedErrContent := `{"width": "test"
<-->
}`
	if rErr.Content != expectedErrContent {
		t.Errorf("Expected error to be '%s', got '%s'", expectedErrContent, rErr.Content)
	}
}

func TestValidAPIKey(t *testing.T) {
	r := test.NewResponse(http.StatusOK, "./fixtures/GetValidAPIKey.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /": r})
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

func TestValidAPIKey_InvalidKey(t *testing.T) {
	r := test.NewResponse(http.StatusForbidden, "./fixtures/GetInvalidAPIKey.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /": r})
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
