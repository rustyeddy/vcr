package main

import (
	"net/http/httptest"
	"testing"
)

// TestHandlers will test a few of the HTTP handlers
func TestHandlers(t *testing.T) {
	ts := httptest.NewServer()

	srvQ := NewHTTPServer(config)

}
