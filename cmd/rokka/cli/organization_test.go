package cli

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/rokka-io/rokka-go/test"
	"github.com/spf13/cobra"
)

func TestCreateOrganization(t *testing.T) {
	org := "test-org"
	email := "testorg@example.org"
	displayName := "Test Org"
	ts := test.NewMockAPI(test.Routes{"PUT /organizations/" + org: test.Response{http.StatusOK, "../../../rokka/fixtures/CreateOrganization.json", nil}})
	defer ts.Close()

	f, err := ioutil.TempFile(os.TempDir(), "stdin")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	log := NewCLILog(false)
	log.StdOut = f
	logger = log

	cl = rokka.NewClient(&rokka.Config{APIAddress: ts.URL, HTTPClient: newHTTPClient(log)})

	organizationCreateCmd.Run(&cobra.Command{}, []string{org, email, displayName})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)

	expected := "Id:             8c0cbde4-ba62-11e7-abc4-cec278b6b50a\nName:           test\nDisplay name:   Dev Environment\nBilling email:  info@example.com\n"
	if string(b) != expected {
		t.Errorf("Expected stdout output to match '%s', got '%s'", expected, b)
	}
}
