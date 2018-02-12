package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

// Response contains a response the testserver responds with
type Response struct {
	StatusCode int
	FileName   string
	Headers    map[string]string
}

// Routes is a map whose key is method + path. If that matches, the assigned response is returned
type Routes map[string]Response

// NewMockAPI starts a httptest server which responds to any request with the given response.
func NewMockAPI(responses Routes) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, ok := responses[r.Method+" "+r.URL.String()]
		if !ok {
			response, ok = responses[r.Method+" "+r.URL.Path]
		}
		if ok {
			file := []byte{}
			var err error
			if response.FileName != "" {
				file, err = ioutil.ReadFile(response.FileName)
				if err != nil {
					log.Fatal(err)
				}
			}

			w.Header().Set("Content-Type", "application/json")
			for key, value := range response.Headers {
				w.Header().Set(key, value)
			}
			if response.StatusCode != 0 {
				w.WriteHeader(response.StatusCode)
			}
			w.Write([]byte(file))
			return
		}
		fmt.Printf("NOT FOUND %s %s\n", r.Method, r.URL)
		w.WriteHeader(http.StatusNotFound)
	}))
}
