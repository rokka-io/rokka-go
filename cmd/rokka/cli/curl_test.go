package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/rokka-io/rokka-go/test"
	"github.com/spf13/cobra"
)

func TestCurl(t *testing.T) {
	fileName := "../../../rokka/fixtures/GetOperations.json"
	headerNameToSet := "Rokka-Headerli"
	headerValueToSet := "is set"
	r := test.NewResponse(http.StatusOK, fileName)
	r.Assertion = func(t *testing.T, r *http.Request) {
		got := r.Header.Get(headerNameToSet)
		if got != headerValueToSet {
			t.Errorf("Expected header '%s: %s', got '%s'", headerNameToSet, headerValueToSet, got)
		}
	}
	ts := test.NewMockAPI(t, test.Routes{"GET /operations": r})
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

	cmd := curlCmd
	if err := cmd.Flags().Set("header", fmt.Sprintf("%s: %s", headerNameToSet, headerValueToSet)); err != nil {
		panic(err)
	}
	defer cmd.Flags().Set("header", "") // reset flag to default value

	cmd.Run(&cobra.Command{}, []string{"/operations"})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	expected, err := ioutil.ReadFile("../../../rokka/fixtures/GetOperations.json")
	if err != nil {
		panic(err)
	}
	if string(b) != string(expected) {
		t.Errorf("Expected stdout output to match '%s', got '%s'", expected, b)
	}
}

func TestCurl_IncludeHeaders(t *testing.T) {
	fileName := "../../../rokka/fixtures/GetOperations.json"
	r := test.NewResponse(http.StatusOK, fileName)
	// override auto-generated Date Header for easy testing.
	r.Headers["Date"] = "Mon, 19 Feb 2018 17:45:52 GMT"
	ts := test.NewMockAPI(t, test.Routes{"GET /operations": r})
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

	cmd := curlCmd
	if err := cmd.Flags().Set("include", "true"); err != nil {
		panic(err)
	}
	defer cmd.Flags().Set("include", "false") // reset flag to default value
	cmd.Run(&cobra.Command{}, []string{"/operations"})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	expected, err := ioutil.ReadFile("../../../rokka/fixtures/GetOperations.json")
	if err != nil {
		panic(err)
	}

	expectedAnnotated := "HTTP/1.1 200 OK\nContent-Type: application/json\nDate: Mon, 19 Feb 2018 17:45:52 GMT\n\n" + string(expected)

	if string(b) != expectedAnnotated {
		t.Errorf("Expected stdout output to match '%s', got '%s'", expectedAnnotated, b)
	}
}

func TestCurl_OverrideMethod(t *testing.T) {
	fileName := "../../../rokka/fixtures/GetOperations.json"
	r := test.NewResponse(http.StatusOK, fileName)
	ts := test.NewMockAPI(t, test.Routes{"PUT /operations": r})
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

	cmd := curlCmd
	if err := cmd.Flags().Set("request", "PUT"); err != nil {
		panic(err)
	}
	defer cmd.Flags().Set("request", "GET") // reset flag to default value
	cmd.Run(&cobra.Command{}, []string{"/operations"})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	expected, err := ioutil.ReadFile("../../../rokka/fixtures/GetOperations.json")
	if err != nil {
		panic(err)
	}
	if string(b) != string(expected) {
		t.Errorf("Expected stdout output to match '%s', got '%s'", expected, b)
	}
}

func TestCurl_SendDataString(t *testing.T) {
	expectedBody := `{"foo": "bar"}`
	r := test.NewResponse(http.StatusOK, "")
	r.Assertion = func(t *testing.T, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if string(b) != expectedBody {
			t.Errorf("Expected body to be '%s', got: '%s'", expectedBody, string(b))
		}
	}
	ts := test.NewMockAPI(t, test.Routes{"PUT /operations": r})
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

	cmd := curlCmd
	if err := cmd.Flags().Set("data", expectedBody); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set("request", "PUT"); err != nil {
		panic(err)
	}
	defer func() {
		// reset flag to default value
		cmd.Flags().Set("data", "")
		cmd.Flags().Set("request", "GET")
	}()

	cmd.Run(&cobra.Command{}, []string{"/operations"})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if string(b) != "" {
		t.Errorf("Expected stdout output to match '%s', got '%s'", "", b)
	}
}

func TestCurl_SendDataFile(t *testing.T) {
	expectedBody := `{"foo": "bar"}`
	inputFile, err := ioutil.TempFile(os.TempDir(), "test-input-file")
	if err != nil {
		panic(err)
	}
	defer os.Remove(inputFile.Name())

	_, err = inputFile.WriteString(expectedBody)
	if err != nil {
		panic(err)
	}

	r := test.NewResponse(http.StatusOK, "")
	r.Assertion = func(t *testing.T, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if string(b) != expectedBody {
			t.Errorf("Expected body to be '%s', got: '%s'", expectedBody, string(b))
		}
	}
	ts := test.NewMockAPI(t, test.Routes{"PUT /operations": r})
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

	cmd := curlCmd
	if err := cmd.Flags().Set("data", "@"+inputFile.Name()); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set("request", "PUT"); err != nil {
		panic(err)
	}
	defer func() {
		// reset flag to default value
		cmd.Flags().Set("data", "")
		cmd.Flags().Set("request", "GET")
	}()

	cmd.Run(&cobra.Command{}, []string{"/operations"})

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if string(b) != "" {
		t.Errorf("Expected stdout output to match '%s', got '%s'", "", b)
	}
}
