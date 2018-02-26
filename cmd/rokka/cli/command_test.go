package cli

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/spf13/cobra"
)

func TestTitleCamelCase(t *testing.T) {
	exp := "TestFooBar"
	act := TitleCamelCase("test-foo_bar")
	if exp != act {
		t.Errorf("Expected '%s', got: '%s'", exp, act)
	}
}

func TestRun_ValidOutput(t *testing.T) {
	org := "test-org"
	fn := func(c *rokka.Client, args []string) (interface{}, error) {
		if args[0] != org {
			t.Errorf("Unexpected value for args[0], got: '%s', expected: '%s'", args[0], org)
		}
		return "Test", nil
	}
	stdErr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		panic(err)
	}
	defer os.Remove(stdErr.Name())

	stdOut, err := ioutil.TempFile(os.TempDir(), "stdout")
	if err != nil {
		panic(err)
	}
	defer os.Remove(stdOut.Name())

	logger = newCLILog(false)
	logger.StdErr = stdErr
	logger.StdOut = stdOut

	cmd := cobra.Command{}

	cFn := run(fn, "{{.}}")
	cFn(&cmd, []string{org})

	stdOut.Seek(0, 0)
	stdErr.Seek(0, 0)

	stdErrOut, err := ioutil.ReadAll(stdErr)
	if err != nil {
		panic(err)
	}
	stdOutOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		panic(err)
	}
	if string(stdErrOut) != "" {
		t.Errorf("Expected empty stdErr, got: '%s'", stdErrOut)
	}
	if string(stdOutOut) != "Test" {
		t.Errorf("Expected 'Test' in stdOut, got: '%s'", stdOutOut)
	}
}

func TestRun_ErrorOutput(t *testing.T) {
	if _, ok := os.LookupEnv("EXECUTE_CRASHING_FUNCTION"); !ok {
		cmd := exec.Command(os.Args[0], "-test.run=TestRun_ErrorOutput")
		cmd.Env = append(os.Environ(), "EXECUTE_CRASHING_FUNCTION=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}

	fn := func(c *rokka.Client, args []string) (interface{}, error) {
		return nil, errors.New("some error for testing")
	}
	stdErr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		panic(err)
	}
	defer os.Remove(stdErr.Name())

	stdOut, err := ioutil.TempFile(os.TempDir(), "stdout")
	if err != nil {
		panic(err)
	}
	defer os.Remove(stdOut.Name())

	logger = newCLILog(false)

	cmd := cobra.Command{}

	cFn := run(fn, "{{.}}")
	cFn(&cmd, []string{})
}
