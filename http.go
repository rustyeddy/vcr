package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// Return an
type HTTPServer struct {
	Addr string

	*httprouter.Router
	Q chan string
}

// NewHTTPServer creates a new HTTP Server
func NewHTTPServer(config *Configuration) (s *HTTPServer) {
	log.Info().
		Str("Addr", config.Addr).
		Str("StaticPath", config.StaticPath).
		Msg("New HTTP Server created")

	// If Q is nil then the server is not running
	s = &HTTPServer{
		Router: httprouter.New(),
		Addr:   config.Addr,
		Q:      nil,
	}

	s.AddHandler("/health", health)
	s.AddHandler("/config", getConfig)

	return s
}

// Start the HTTP server, give the caller a channel back that will
// allow the caller to communicate with this server
func (s *HTTPServer) Start(wg *sync.WaitGroup) chan string {
	defer wg.Done()

	log.Info().Msg("HTTP Server is starting")

	s.Q = make(chan string)
	go func() {
		// Blocks unless something goes wrong
		log.Info().Msg("\tgo routine http listen ...")
		if err := http.ListenAndServe(s.Addr, s.Router); err != nil {
			log.Error().Msg("Errored on the listen.")
		}
	}()
	log.Info().Msg("non-go returning s.Q")
	return s.Q
}

// AddHandler
func (s *HTTPServer) AddHandler(path string, f httprouter.Handle) {
	s.GET(path, f)
}

func health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := map[string]string{
		"health": "ok",
	}
	json.NewEncoder(w).Encode(resp)
}

// healthCheckHndl
func getConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(config)
}
