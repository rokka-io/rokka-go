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
	r := test.NewResponse(http.StatusOK, "../../../rokka/fixtures/CreateOrganization.json")
	ts := test.NewMockAPI(t, test.Routes{"PUT /organizations/" + org: r})
	defer ts.Close()

	f, err := ioutil.TempFile(os.TempDir(), "stdin")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	log := newCLILog(false)
	log.StdOut = f
	logger = log

	rokkaClient = rokka.NewClient(&rokka.Config{APIAddress: ts.URL, HTTPClient: newHTTPClient(log)})

	organizationCreateCmd.Run(&cobra.Command{}, []string{org, email, displayName})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	expected := "Id:             8c0cbde4-ba62-11e7-abc4-cec278b6b50a\nName:           test\nDisplay name:   Dev Environment\nBilling email:  info@example.com\n"
	if string(b) != expected {
		t.Errorf("Expected stdout output to match '%s', got '%s'", expected, b)
	}
}
