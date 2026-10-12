package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessReturnsBody(t *testing.T) {
	got, err := fetchAll(serve(t, body("hello\n")))
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if got != "hello\n" {
		t.Errorf("body = %q, want %q", got, "hello\n")
	}
}

func TestNon2xxStatusIsError(t *testing.T) {
	for _, code := range []int{http.StatusNotFound, http.StatusInternalServerError} {
		if _, err := fetchAll(serve(t, status(code))); err == nil {
			t.Errorf("status %d: fetch succeeded, want error", code)
		}
	}
}

func fetchAll(url string) (string, error) {
	r, err := fetch(context.Background(), url)
	if err != nil {
		return "", err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	return string(b), err
}

func serve(t *testing.T, h http.HandlerFunc) string {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv.URL
}

func body(s string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}
}

func status(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(code), code)
	}
}
