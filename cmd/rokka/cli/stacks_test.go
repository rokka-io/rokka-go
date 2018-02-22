package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/rokka-io/rokka-go/rokka"
	"github.com/rokka-io/rokka-go/test"
)

func TestCreateStack_Piped(t *testing.T) {
	org := "test-org"
	name := "test-name"
	r := test.NewResponse(http.StatusOK, "../../../rokka/fixtures/CreateStack.json")
	ts := test.NewMockAPI(t, test.Routes{"PUT /stacks/" + org + "/" + name: r})
	defer ts.Close()
	c := rokka.NewClient(&rokka.Config{APIAddress: ts.URL})

	f, err := ioutil.TempFile(os.TempDir(), "stdin")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	// override the stdin createStack uses to our temp file.
	stdin = f

	ops := rokka.Operations{
		rokka.AlphaOperation{
			Mode: rokka.StrPtr("mask"),
		},
		rokka.ResizeOperation{
			Height: rokka.IntPtr(100),
			Width:  rokka.IntPtr(120),
		},
	}
	options := rokka.StackOptions{
		"basestack": "test-basestack",
	}
	expressions := make([]rokka.Expression, 1)
	expressions[0] = rokka.Expression{
		Expression: "options.dpr >= 2",
		Overrides: map[string]interface{}{
			"options": map[string]interface{}{
				"jpq.quality": 60,
			},
		},
	}

	stack := rokka.CreateStackRequest{
		Operations:  ops,
		Options:     options,
		Expressions: expressions,
	}

	b, err := json.Marshal(stack)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
	// reset read/write pointer to start
	f.Seek(0, 0)

	res, err := createStack(c, []string{org, name})
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v", res)
}

func TestCreateStack_CLI(t *testing.T) {
	org := "test-org"
	name := "test-name"
	r := test.NewResponse(http.StatusOK, "../../../rokka/fixtures/GetOperations.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /operations": r})
	defer ts.Close()
	c := rokka.NewClient(&rokka.Config{APIAddress: ts.URL})

	f, err := ioutil.TempFile(os.TempDir(), "stdin")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	stdin = f

	inputs := []string{
		"a",     // add alpha operation
		"mask",  // mode option
		"trim",  // add trim operation
		"1.0",   // fuzzy option
		"resiz", // add resize operation
		"1",     // height option
		"box",   // mode option
		"false", // upscale option
		"true",  // upscale_dpr option
		"asdf",  // width option - wrong input
		"y",     // yes, want to retry
		"2",     // width option - correct input
		"ro",    // rotate operation
		"",      // angle (required) - empty input
		"",      // background_color - empty input
		"",      // background_opacity - empty input
		"y",     // validation failed because angle is required, yes want to retry
		"",      // angle (required) - empty input
		"",      // background_color - empty input
		"",      // background_opacity - empty input
		"n",     // validation failed because angle is required, no don't want to retry
		"q",     // quit adding new stuff
		"",      // newline in the end is required
	}
	input := strings.Join(inputs, "\n")

	_, err = f.Write([]byte(input))
	if err != nil {
		panic(err)
	}
	f.Seek(0, 0)

	req := rokka.CreateStackRequest{}

	// unfortunately can't use `createStack` because with a already written stdin-fake `createStack` will guess it's piped JSON data.
	err = cliCreateStack(c, name, org, &req)
	if err != nil {
		t.Error(err)
	}

	exp := 3
	if len(req.Operations) != exp {
		t.Fatalf("Expected %d operation to be added, got %d", exp, len(req.Operations))
	}
	opAlpha, ok := req.Operations[0].(*rokka.AlphaOperation)
	if !ok {
		t.Fatalf("Expected operation to be of type *rokka.AlphaOperation, got: %T", req.Operations[0])
	}
	if *opAlpha.Mode != "mask" {
		t.Fatalf("Expected Mode to be %s, got %s", "mask", *opAlpha.Mode)
	}

	opTrim, ok := req.Operations[1].(*rokka.TrimOperation)
	if !ok {
		t.Fatalf("Expected operation to be of type *rokka.TrimOperation, got: %T", req.Operations[1])
	}
	if *opTrim.Fuzzy != 1.0 {
		t.Fatalf("Expected Fuzzy to be %f, got %f", 1.0, *opTrim.Fuzzy)
	}

	opResize, ok := req.Operations[2].(*rokka.ResizeOperation)
	if !ok {
		t.Fatalf("Expected operation to be of type *rokka.ResizeOperation, got: %T", req.Operations[2])
	}
	if *opResize.Height != 1 {
		t.Fatalf("Expected Height to be %d, got %d", 1, *opResize.Height)
	}
	if *opResize.Mode != "box" {
		t.Fatalf("Expected Mode to be %s, got %s", "mode", *opResize.Mode)
	}
	if *opResize.Upscale != false {
		t.Fatalf("Expected Upscale to be %t, got %t", false, *opResize.Upscale)
	}
	if *opResize.UpscaleDpr != true {
		t.Fatalf("Expected UpscaleDpr to be %t, got %t", true, *opResize.UpscaleDpr)
	}
	if *opResize.Width != 2 {
		t.Fatalf("Expected Width to be %d, got %d", 2, *opResize.Width)
	}

	if len(req.Expressions) != 0 {
		t.Fatalf("Expected %d expressions to be added, got %d", 0, len(req.Expressions))
	}
	if len(req.Options) != 0 {
		t.Fatalf("Expected %d options to be added, got %d", 0, len(req.Options))
	}
}
