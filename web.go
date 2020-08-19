package main

import (
	"encoding/json"
	"log"
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
func web() {
	log.Println("New HTTP Server %s created", config.Addr)

	http.HandleFunc("/api/health", health)
	http.HandleFunc("/api/video", getVideo)
	http.HandleFunc("/api/video/play", playVideo)
	http.HandleFunc("/api/video/pause", pauseVideo)

	go http.ListenAndServe(config.Addr, nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(successResponse)
}

// getMessanger
func getVideo(w http.ResponseWriter, r *http.Request) {
	status := &VideoPlayerStatus{}
	if vid := GetVideoPlayer(); vid != nil {
		status = vid.Status()
	}
	json.NewEncoder(w).Encode(status)
}

// getMessanger
func playVideo(w http.ResponseWriter, r *http.Request) {
	if vid := GetVideoPlayer(); vid != nil {
		cmdQ <- NewTLV(CMDPlay, 2)
	}
	json.NewEncoder(w).Encode(successResponse)
}

// getMessanger
func pauseVideo(w http.ResponseWriter, r *http.Request) {
	if v := GetVideoPlayer(); v != nil {
		v.Pause()
	}
	json.NewEncoder(w).Encode(successResponse)
}
