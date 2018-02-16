package rokka

import (
	"net/http"
	"testing"

	"github.com/rokka-io/rokka-go/test"
)

func TestListStacks(t *testing.T) {
	org := "test-org"
	ts := test.NewMockAPI(test.Routes{"GET /stacks/" + org: test.Response{http.StatusOK, "./fixtures/ListStacks.json", nil}})
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
	ts := test.NewMockAPI(test.Routes{"PUT /stacks/" + org + "/" + name: test.Response{http.StatusOK, "./fixtures/CreateStack.json", nil}})
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
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestDeleteStack(t *testing.T) {
	org := "test-org"
	name := "test-stack"
	ts := test.NewMockAPI(test.Routes{"DELETE /stacks/" + org + "/" + name: test.Response{http.StatusNoContent, "", nil}})
	defer ts.Close()

	c := NewClient(&Config{APIAddress: ts.URL})

	err := c.DeleteStack(org, name)
	if err != nil {
		t.Error(err)
	}
}
