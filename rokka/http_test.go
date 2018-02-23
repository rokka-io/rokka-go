package rokka

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rokka-io/rokka-go/test"
)

func TestRetry_TooManyRequests(t *testing.T) {
	r := test.NewResponse(http.StatusTooManyRequests, "")
	r.Headers["Retry-After"] = "0.1"
	ts := test.NewMockAPI(t, test.Routes{"GET /": r})
	defer ts.Close()

	c := NewClient(&Config{
		APIAddress:         ts.URL,
		APIKey:             "test",
		RetryingHTTPClient: NewRetryingHTTPClient(DefaultConfig().HTTPClient, 10, 1),
	})
	retryingClient := c.AutoRetry()

	req, err := retryingClient.NewRequest(http.MethodGet, "/", nil, nil)
	if err != nil {
		panic(err)
	}

	err = retryingClient.Call(req, nil, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	maxRetriesErr, ok := err.(ErrMaxRetriesReached)
	if !ok {
		t.Fatalf("Expected error of type '%T', got '%T'", ErrMaxRetriesReached{}, err)
	}
	if maxRetriesErr.LastError != nil {
		t.Errorf("Expected no LastError, got '%v'", maxRetriesErr.LastError)
	}
}

func TestRetry_BadGateway(t *testing.T) {
	r := test.NewResponse(http.StatusBadGateway, "")
	ts := test.NewMockAPI(t, test.Routes{"GET /": r})
	defer ts.Close()

	c := NewClient(&Config{
		APIAddress:         ts.URL,
		APIKey:             "test",
		RetryingHTTPClient: NewRetryingHTTPClient(DefaultConfig().HTTPClient, 10, 1),
	})
	retryingClient := c.AutoRetry()

	req, err := retryingClient.NewRequest(http.MethodGet, "/", nil, nil)
	if err != nil {
		panic(err)
	}

	err = retryingClient.Call(req, nil, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	maxRetriesErr, ok := err.(ErrMaxRetriesReached)
	if !ok {
		t.Fatalf("Expected error of type '%T', got '%T'", ErrMaxRetriesReached{}, err)
	}
	if maxRetriesErr.LastError != nil {
		t.Errorf("Expected no LastError, got '%v'", maxRetriesErr.LastError)
	}
}

func TestRetry_NoConnection(t *testing.T) {
	hClient := &http.Client{
		Timeout: 1 * time.Microsecond,
	}
	c := NewClient(&Config{
		APIAddress:         "http://localhost:0",
		APIKey:             "test",
		HTTPClient:         hClient,
		RetryingHTTPClient: NewRetryingHTTPClient(hClient, 10, 1),
	})
	retryingClient := c.AutoRetry()

	req, err := retryingClient.NewRequest(http.MethodGet, "/", nil, nil)
	if err != nil {
		panic(err)
	}

	err = retryingClient.Call(req, nil, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	maxRetriesErr, ok := err.(ErrMaxRetriesReached)
	if !ok {
		t.Fatalf("Expected error of type '%T', got '%T'", ErrMaxRetriesReached{}, err)
	}
	if maxRetriesErr.LastError == nil {
		t.Fatalf("Expected LastError, got nil")
	}
	if !strings.Contains(maxRetriesErr.LastError.Error(), "Get http://localhost:0/:") {
		t.Errorf("Expected LastError, got: '%s'", maxRetriesErr.LastError.Error())
	}
}

func TestCalculateBackoff(t *testing.T) {
	table := []struct {
		retries  int
		max      int
		expected int
	}{
		{
			retries:  0,
			max:      10000,
			expected: 0,
		},
		{
			retries:  1,
			max:      10000,
			expected: 200,
		},
		{
			retries:  2,
			max:      10000,
			expected: 440,
		},
		{
			retries:  5,
			max:      10000,
			expected: 1488,
		},
		{
			retries:  10,
			max:      10000,
			expected: 5192,
		},
		{
			retries:  13,
			max:      10000,
			expected: 9699,
		},
		{
			retries:  14,
			max:      10000,
			expected: 10000,
		},
	}

	for _, v := range table {
		t.Run(fmt.Sprintf("retries: '%d' max: '%d'", v.retries, v.max), func(t *testing.T) {
			actual := calculateBackoff(v.retries, v.max)
			if v.expected != actual {
				t.Errorf("Expected '%d' got '%d'", v.expected, actual)
			}
		})
	}
}
