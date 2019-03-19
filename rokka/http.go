package rokka

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RetryingHTTPClient implements HTTPRequester and wraps another HTTPRequester.
type RetryingHTTPClient struct {
	c          HTTPRequester
	maxRetries int
	maxDelay   int
}

// NewRetryingHTTPClient wraps an HTTPRequester in order to do automatic retries.
//
// These are done in the following cases:
//  - a HTTP, Network or Transport error occurred
//  - a status code of 429 (too many requests), 502 (bad gateway), 503 (service unavailable), or 504 (gateway timeout) has been received.
func NewRetryingHTTPClient(c HTTPRequester, maxRetries int, maxDelay int) *RetryingHTTPClient {
	return &RetryingHTTPClient{
		c:          c,
		maxRetries: maxRetries,
		maxDelay:   maxDelay,
	}
}

// ErrMaxRetriesReached is returned if after maxRetries number of retries the request still fails.
type ErrMaxRetriesReached struct {
	LastError error
}

func (e ErrMaxRetriesReached) Error() string {
	return fmt.Sprintf("rokka: max retries reached. Last error: %s", e.LastError)
}

// Do executes an HTTP request and retries in case the response status code is one 429, 502, 503, 504 or if error is set.
func (hc *RetryingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < hc.maxRetries; i++ {
		resp, err = hc.c.Do(req)

		if !shouldRetry(resp, err) {
			return resp, err
		}

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := retryAfterSeconds(resp.Header.Get("Retry-After"), float64(hc.maxDelay))
			if retryAfter != 0.0 {
				time.Sleep(time.Duration(retryAfter) * time.Second)
				continue
			}
		}
		delay := calculateBackoff(i, hc.maxDelay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	return resp, ErrMaxRetriesReached{LastError: err}
}

// shouldRetry determines in which cases a retry is tried.
// Retries if there was any error (http error, network error, transport error etc.).
// Also retries on codes 429 (too many requests), 502 (bad gateway), 503 (service unavailable), 504 (gateway timeout).
func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	sc := resp.StatusCode
	if sc == http.StatusTooManyRequests ||
		sc == http.StatusBadGateway ||
		sc == http.StatusServiceUnavailable ||
		sc == http.StatusGatewayTimeout {
		return true
	}
	return false
}

// retryAfterSeconds calculates the delay to wait in seconds based on the Retry-After header.
// It takes into account both varieties of the header: number of seconds or a Date.
// To account for possible issues with the Retry-After header, it doesn't wait longer than `max` seconds.
// Will return 0.0 if retryAfter is an empty string or in a format not understood.
func retryAfterSeconds(retryAfter string, max float64) float64 {
	if retryAfter == "" {
		return 0.0
	}

	delay := 0.0
	retryAfterSeconds, err := strconv.Atoi(retryAfter)
	if err == nil {
		delay = float64(retryAfterSeconds)
	} else {
		retryAfterDate, err := http.ParseTime(retryAfter)
		if err == nil {
			delay = time.Until(retryAfterDate).Seconds()
		}
	}

	if max > delay {
		return delay
	}
	return max
}

// calculateBackoff calculates an exponential backoff time in miliseconds with a smaller exponential factor than
// the usual `2` in order not to wait too long.
// Example wait times can be seen in the accompanying test.
// If the calculated wait time exceeds max, will always return max.
func calculateBackoff(retries int, max int) int {
	backoffFactor := math.Pow(1.2, float64(retries)) - 1
	delay := int(round(backoffFactor * 1000))
	if max > delay {
		return delay
	}
	return max
}
