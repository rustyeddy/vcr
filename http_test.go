package main

import (
	"net/http/httptest"
	"sync"
	"testing"
)

// TestHandlers will test a few of the HTTP handlers
func TestHandlers(t *testing.T) {
	var wg sync.WaitGroup
	srv := NewHTTPServer(&config)
	req := httptest.NewRequest("GET", "http://1.2.4.3", nil)
	w := httptest.NewRecorder()
	wg.Add(1)
	srv.Start(&wg)

	Health(w, req, nil)

}
