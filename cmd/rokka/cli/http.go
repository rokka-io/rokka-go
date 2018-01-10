package cli

import (
	"net/http"
	"net/http/httputil"
)

// httpClient implements rokka.HTTPRequester
type httpClient struct {
	c   *http.Client
	log *Log
}

func newHTTPClient(log *Log) *httpClient {
	return &httpClient{
		c:   &http.Client{},
		log: log,
	}
}

// Do executes an HTTP request. If verbose is set this will log the request and response using httputil.
func (hc *httpClient) Do(req *http.Request) (*http.Response, error) {
	if hc.log.Verbose {
		dump, err := httputil.DumpRequest(req, req.Header.Get("Content-Type") == "application/json")
		if err != nil {
			hc.log.Errorf("Unable to dump request: %s\n", err)
		} else {
			hc.log.Errorf("%s\n", dump)
		}
	}

	resp, err := hc.c.Do(req)
	if err == nil && hc.log.Verbose {
		dump, err := httputil.DumpResponse(resp, resp.Header.Get("Content-Type") == "application/json")
		if err != nil {
			hc.log.Errorf("Unable to dump request: %s\n", err)
		} else {
			hc.log.Errorf("%s\n", dump)
		}
	}

	return resp, err
}
