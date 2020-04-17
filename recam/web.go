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
	successResponse = map[string]string{"success": "true"}
	errorResponse   = map[string]string{"success": "false"}
)

// NewWebServer creates a new HTTP Server
func NewWebServer(config *Configuration) (s *WebServer) {
	log.Info().
		Str("Addr", config.Addr).
		Str("StaticPath", config.StaticPath).
		Msg("New HTTP Server created")

	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", r.Header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		log.Info().Msg("A REST Request has been had")
		w.WriteHeader(http.StatusNoContent)
	})

	// If Q is nil then the server is not running
	s = &WebServer{
		Router: router,
		Addr:   config.Addr,
	}

	s.AddHandler("/health", health)
	s.AddHandler("/config", getConfig)
	s.AddHandler("/messanger", getMessanger)
	s.AddHandler("/video", getVideo)
	s.AddHandler("/video/play", playVideo)
	s.AddHandler("/video/pause", pauseVideo)
	return s
}

// Start the HTTP server, give the caller a channel back that will
// allow the caller to communicate with this server
func (s *WebServer) Start() {
	log.Info().Msg("HTTP Server is starting")
	go func() {
		// Blocks unless something goes wrong
		log.Info().
			Str("addr", s.Addr).
			Msg("\tgo routine http listen ...")
		if err := http.ListenAndServe(s.Addr, s.Router); err != nil {
			log.Error().Msg("Errored on the listen.")
		}
	}()
}

// AddHandler
func (s *WebServer) AddHandler(path string, f httprouter.Handle) {
	s.GET(path, f)
}

func health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(successResponse)
}

// healthCheckHndl
func getConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(config)
}

// getMessanger
func getMessanger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status *MessangerStatus
	if msg != nil {
		status = msg.GetStatus()
	} else {
		// serve up the null entry
		status = &MessangerStatus{
			Broker: "DISCONNECTED",
		}
	}
	json.NewEncoder(w).Encode(status)
}

// getMessanger
func getVideo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := &VideoPlayerStatus{}
	if vid != nil {
		status = vid.Status()
	}
	json.NewEncoder(w).Encode(status)
}

// getMessanger
func playVideo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if vid != nil {
		vidQ <- "play"
	}
	json.NewEncoder(w).Encode(successResponse)
}

// getMessanger
func pauseVideo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if vid != nil {
		vidQ <- "pause"
	}
	json.NewEncoder(w).Encode(successResponse)
}
