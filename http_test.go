package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

func init() {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if false {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

}

// TestHandlers will test a few of the HTTP handlers
func TestHealth(t *testing.T) {
	var wg sync.WaitGroup
	srv := NewHTTPServer(&config)
	req := httptest.NewRequest("GET", "http://1.2.4.3", nil)
	w := httptest.NewRecorder()
	wg.Add(1)
	srv.Start(&wg)

	health(w, req, nil)
	if w.Code != 200 {
		t.Errorf("health check failed expected (%d) got (%d)  ", 200, w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("health check body failed error: %d", err)
	}

	// TODO unravel the reponse and check the values
	if !strings.Contains(string(body), "ok") {
		t.Errorf("health check body is wrong got %s", body)
	}
}

// TestHandlers will test a few of the HTTP handlers
func TestConfig(t *testing.T) {
	var wg sync.WaitGroup
	srv := NewHTTPServer(&config)
	req := httptest.NewRequest("GET", "http://1.2.4.3", nil)
	w := httptest.NewRecorder()
	wg.Add(1)
	srv.Start(&wg)

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
