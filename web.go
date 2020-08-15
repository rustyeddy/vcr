package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// Return an
type WebServer struct {
	Addr string
	*httprouter.Router
}

var (
	server          *WebServer
	successResponse = map[string]string{"success": "true"}
	errorResponse   = map[string]string{"success": "false"}
)

// NewWebServer creates a new HTTP Server
func NewWebServer() (s *WebServer) {
	log.Info().
		Str("Addr", config.Addr).
		Str("State", "created").
		Msg("New HTTP Server created")

	http.HandleFunc("/api/health", health)
	http.HandleFunc("/api/messanger", getMessanger)
	http.HandleFunc("/api/video", getVideo)
	http.HandleFunc("/api/video/play", playVideo)
	http.HandleFunc("/api/video/pause", pauseVideo)
	return s
}

// Start the HTTP server, give the caller a channel back that will
// allow the caller to communicate with this server
func (s *WebServer) Start(cmdQ chan TLV) chan TLV {
	log.Info().Str("addr", s.Addr).Str("state", "start").Msg("Starting the HTTP Server...")

	q := make(chan TLV)
	go func() {
		var cmd TLV
		for {
			src := ""
			log.Info().Msg("\tWeb Server listening for internal communication.")
			select {
			case cmd = <-q:
				src = "webQ"
			}
			log.Info().Str("cmd", cmd.Str()).Str("src", src).Msg("Incoming command")
		}
	}()

	// The gofunc the listen code
	go func() {
		// Blocks unless something goes wrong
		log.Info().
			Str("addr", s.Addr).
			Msg("\tgo routine http listen ...")
		if err := http.ListenAndServe(s.Addr, s.Router); err != nil {
			log.Error().Msg("Errored on the listen.")
		}
	}()
	return q
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(successResponse)
}

// healthCheckHndl
func getConfig(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Config)
}

// getMessanger
func getMessanger(w http.ResponseWriter, r *http.Request) {
	var status *MessangerStatus
	if m := GetMessanger(); m != nil {
		status = m.GetStatus()
	} else {
		// serve up the null entry
		status = &MessangerStatus{
			Broker: "DISCONNECTED",
		}
	}
	json.NewEncoder(w).Encode(status)
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
