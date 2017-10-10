package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
)

func NewMockAPI(file string, statusCode int) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		w.WriteHeader(statusCode)
	}))

	return ts
}