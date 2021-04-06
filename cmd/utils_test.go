package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func givenHTTPServer(t *testing.T, h http.Handler) *httptest.Server {
	server := httptest.NewServer(h)
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

func givenStaticContentHandler(content string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", content)
	})
}
