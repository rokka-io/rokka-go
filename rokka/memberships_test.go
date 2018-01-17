package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestCreateMembership(t *testing.T) {
	ts := test.NewMockAPI("./fixtures/CreateMembership.json", http.StatusCreated)
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.CreateMembership("test", "test@example.org", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Successfully created membership")
}
