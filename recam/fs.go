package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// Return an
type FileServer struct {
	Addr string

	*httprouter.Router
	Q chan string
}

func NewFileServer(config *Configuration) (fs *FileServer) {
	router := httprouter.New()
	/*
		router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Access-Control-Request-Method") != "" {
				// Set CORS headers
				header := w.Header()
				header.Set("Access-Control-Allow-Methods", r.Header.Get("Allow"))
				header.Set("Access-Control-Allow-Origin", "*")
			}

			// Adjust status code to 204
			w.WriteHeader(http.StatusNoContent)
		})*/
	fs = &FileServer{
		Router: router,
		Addr:   ":8000",
	}
	fs.ServeFiles("/*filepath", http.Dir(config.StaticPath))
	return fs
}

// Start the HTTP server, give the caller a channel back that will
// allow the caller to communicate with this server
func (fs *FileServer) Start() {
	log.Info().Msg("HTTP Server is starting")
	go func() {

		// Blocks unless something goes wrong
		log.Info().
			Str("addr", fs.Addr).
			Msg("\tgo routine http listen ...")

		if err := http.ListenAndServe(fs.Addr, fs.Router); err != nil {
			log.Error().Msg("Errored on the listen.")
		}
	}()
}
