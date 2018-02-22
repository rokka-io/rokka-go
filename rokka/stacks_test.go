package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListStacks(t *testing.T) {
	org := "test-org"
	r := test.NewResponse(http.StatusOK, "./fixtures/ListStacks.json")
	ts := test.NewMockAPI(t, test.Routes{"GET /stacks/" + org: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	res, err := c.ListStacks(org)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestCreateStacks(t *testing.T) {
	org := "test-org"
	name := "test-stack"
	r := test.NewResponse(http.StatusOK, "./fixtures/CreateStack.json")
	ts := test.NewMockAPI(t, test.Routes{"PUT /stacks/" + org + "/" + name: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	ops := Operations{
		AlphaOperation{
			Mode: StrPtr("mask"),
		},
		ResizeOperation{
			Height: IntPtr(100),
			Width:  IntPtr(120),
		},
	}
	options := StackOptions{
		"basestack": "test-basestack",
	}
	expressions := make([]Expression, 1)
	expressions[0] = Expression{
		Expression: "options.dpr >= 2",
		Overrides: map[string]interface{}{
			"options": map[string]interface{}{
				"jpq.quality": 60,
			},
		},
	}

	res, err := c.CreateStack(org, name, CreateStackRequest{
		Operations:  ops,
		Options:     options,
		Expressions: expressions,
	}, false)
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestDeleteStack(t *testing.T) {
	org := "test-org"
	name := "test-stack"
	r := test.NewResponse(http.StatusNoContent, "")
	ts := test.NewMockAPI(t, test.Routes{"DELETE /stacks/" + org + "/" + name: r})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.DeleteStack(org, name)
	if err != nil {
		t.Error(err)
	}
}
