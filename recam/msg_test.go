package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandlers will test a few of the HTTP handlers
func TestMessanger(t *testing.T) {
	srv := NewMessanger(&config)
	req := httptest.NewRequest("GET", "http://localhost:8888/messanger", nil)
	w := httptest.NewRecorder()
	srv.Start()

	getConfig(w, req, nil)
	if w.Code != 200 {
		t.Errorf("health check failed expected (%d) got (%d)  ", 200, w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("health check body failed error: %d", err)
	}

	// TODO unravel the reponse and check the values
	if !strings.Contains(string(body), "index.html") {
		t.Errorf("health check body is wrong got %s", body)
	}
}
