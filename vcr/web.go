package main

import (
	"log"
	"sync"
	"encoding/json"
	"net/http"
)

// Return an
type WebServer struct {
	Addr string
}

var (
	server          *WebServer
	successResponse = map[string]string{"success": "true"}
	errorResponse   = map[string]string{"success": "false"}
)

// NewWebServer creates a new HTTP Server
func web(wg sync.WaitGroup) {
	defer wg.Done()
	
	log.Printf("New HTTP Server %s created", config.Addr)

	http.HandleFunc("/api/health", health)
	http.ListenAndServe(config.Addr, nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(successResponse)
}
