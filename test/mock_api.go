package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

func NewMockAPI(file string, statusCode int) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}))

	return ts
}
