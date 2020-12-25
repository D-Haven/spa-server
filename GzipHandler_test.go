package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapperAcceptsGzip(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	rr := httptest.NewRecorder()
	handler := GzipHandler(func(w http.ResponseWriter, r *http.Request) {})

	handler.ServeHTTP(rr, req)

	if rr.Result().Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Content was not compressed with gzip")
	}
}

func TestWrapperDoesNotCompressIfRequestDoesNotAcceptIt(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := GzipHandler(func(w http.ResponseWriter, r *http.Request) {})

	handler.ServeHTTP(rr, req)

	if rr.Result().Header.Get("Content-Encoding") == "gzip" {
		t.Errorf("Content was not compressed with gzip")
	}
}
