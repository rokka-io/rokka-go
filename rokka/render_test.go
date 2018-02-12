package rokka

import (
	"testing"
)

func TestGetURLWithoutStackOperations(t *testing.T) {
	expectedURL := "https://test.rokka.io/dynamic/noop/8bbff49a384a4682fd05144ffe77a84f29f112ff.jpg"
	operations := []Operation{}

	c := NewClient(&Config{})
	url, err := c.GetURL("test", "8bbff49a384a4682fd05144ffe77a84f29f112ff", "jpg", operations)
	if err != nil {
		t.Error(err)
	}
	if url != expectedURL {
		t.Errorf("Result doesn't match expected value. Got: \"%s\"; Expected: \"%s\"", url, expectedURL)
	}
}

func TestGetURLWithValidStackOperations(t *testing.T) {
	expectedURL := "https://test.rokka.io/dynamic/composition-height-200-mode-test-width-100--trim--primitive-count-10/8bbff49a384a4682fd05144ffe77a84f29f112ff.png"
	operations := []Operation{
		CompositionOperation{Mode: StrPtr("test"), Width: IntPtr(100), Height: IntPtr(200)},
		TrimOperation{},
		PrimitiveOperation{Count: IntPtr(10)},
	}

	c := NewClient(&Config{})
	url, err := c.GetURL("test", "8bbff49a384a4682fd05144ffe77a84f29f112ff", "png", operations)
	if err != nil {
		t.Error(err)
	}
	if url != expectedURL {
		t.Errorf("Result doesn't match expected value. Got: \"%s\"; Expected: \"%s\"", url, expectedURL)
	}
}

func TestGetURLWithCustomImageHost(t *testing.T) {
	expectedURL := "https://test.example.com/dynamic/noop/8bbff49a384a4682fd05144ffe77a84f29f112ff.png"
	operations := []Operation{}

	c := NewClient(&Config{ImageHost: "https://{{organization}}.example.com"})
	url, err := c.GetURL("test", "8bbff49a384a4682fd05144ffe77a84f29f112ff", "png", operations)
	if err != nil {
		t.Error(err)
	}
	if url != expectedURL {
		t.Errorf("Result doesn't match expected value. Got: \"%s\"; Expected: \"%s\"", url, expectedURL)
	}
}

func TestGetURLWithInvalidStackOperation(t *testing.T) {
	operations := []Operation{
		CompositionOperation{},
	}

	c := NewClient(&Config{})
	_, err := c.GetURL("test", "8bbff49a384a4682fd05144ffe77a84f29f112ff", "png", operations)
	if err == nil {
		t.Error("Error expected")
	}
}

func TestGetURLForStackWithValidStackOperations(t *testing.T) {
	expectedURL := "https://test.rokka.io/stack-name/composition-height-200-mode-test-width-100--trim--primitive-count-10/8bbff49a384a4682fd05144ffe77a84f29f112ff.png"
	operations := []Operation{
		CompositionOperation{Mode: StrPtr("test"), Width: IntPtr(100), Height: IntPtr(200)},
		TrimOperation{},
		PrimitiveOperation{Count: IntPtr(10)},
	}

	c := NewClient(&Config{})
	url, err := c.GetURLForStack("test", "8bbff49a384a4682fd05144ffe77a84f29f112ff", "png", "stack-name", operations)
	if err != nil {
		t.Error(err)
	}
	if url != expectedURL {
		t.Errorf("Result doesn't match expected value. Got: \"%s\"; Expected: \"%s\"", url, expectedURL)
	}
}
