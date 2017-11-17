package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/rokka-io/rokka-go/cli"
)

// HTTPClient implements rokka.HTTPRequester
type HTTPClient struct {
	c   *http.Client
	log *cli.Log
}

func NewHTTPClient(log *cli.Log) *HTTPClient {
	return &HTTPClient{
		c:   &http.Client{},
		log: log,
	}
}

func (hc *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if hc.log.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			hc.log.Errorf("Unable to dump request: %s\n", err)
		} else {
			hc.log.Errorf("%s\n", dump)
		}
	}

	resp, err := hc.c.Do(req)
	if err == nil && hc.log.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			hc.log.Errorf("Unable to dump request: %s\n", err)
		} else {
			hc.log.Errorf("%s\n", dump)
		}
	}

	return resp, err
}
