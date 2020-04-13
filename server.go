package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	// HTTP Server and router handle all HTML and Static page requests
	// including all CSS and JavaScript files.
	*http.Server
	*mux.Router
}

var (
	server *Server
)

// NewServer will create the Server struct, fill in the address, create
// a router, register the route handlers then return
func StartHTTP(wg *sync.WaitGroup, config *Configuration) {
	defer wg.Done()

	srv := &Server{}
	srv.Router = mux.NewRouter()
	srv.Server = &http.Server{
		Handler: srv.Router,
		Addr:    config.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.AddRoute("/ws", wsUpgradeHndl)
	srv.AddRoute("/api/health", healthCheckHndl)
	srv.AddRoute("/api/config", getConfigHndl)
	srv.AddRoute("/api/camera/status", getCameraHndl)
	srv.AddRoute("/api/camera/play", setPlayHndl)
	srv.AddRoute("/api/camera/pause", setPauseHndl)
	srv.AddRoute("/api/camera/snap", setSnapHndl)
	srv.AddRoute("/api/pipelines", getPipelinesHndl)

	// app = spaHandler{StaticPath: "pub", IndexPath: "index.html"}
	srv.Router.PathPrefix("/").Handler(srv)
	log.Print("New Server created")

	// Startup HTTP server
	l.WithField("addr", srv.Addr).Info("Starting HTTP Server ...")
	err := srv.Server.ListenAndServe()
	l.Fatal(err.Error())
}

// AddRoute allows us to dynamically add routes at runtime via a plugin
func (srv *Server) AddRoute(path string, handlr func(http.ResponseWriter, *http.Request)) {
	log.Info().Str("path", path).Msg("adding handler")
	srv.Router.HandleFunc(path, handlr)
}

// StartServer is a meta starter, it starts the HTTP for the SPA and
// the MJPEG server. TODO: should /mjpeg and index.html be on the same
// server (port)?
func (srv *Server) Start(wg *sync.WaitGroup) {
	defer wg.Done()

}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		log.Error().Str("Status", "Bad Request").Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(config.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(config.StaticPath, config.IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(config.StaticPath)).ServeHTTP(w, r)
}

//
// ======================== Handlers =============================
//

// healthCheckHndl
func healthCheckHndl(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// healthCheckHndl
func getConfigHndl(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(config)
}

func setConfigValsHndl(w http.ResponseWriter, r *http.Request) {
	// Assume we accept a config "json" of replacement strings
	log.Fatal().Msg("TODO config vals")
}

func setPlayHndl(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	if !video.Recording {
		// block otherwise because the video player creation and play
		// loop are in this same function, which maybe should be
		// separated?
		go video.StartVideo()
	}
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func setPauseHndl(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	if video.Recording {
		video.StopVideo()
	}
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func setSnapHndl(w http.ResponseWriter, r *http.Request) {
	video.SnapRequest = true
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// getCameraStatus will return the real-time information from the camera
func getCameraHndl(w http.ResponseWriter, r *http.Request) {
	var s CameraStatus
	if video != nil {
		s.Name = video.Name
		s.Addr = video.Addr
		s.Status = "Paused"
		if video.Recording {
			s.Status = "Playing"
		}
		if video.VideoPipeline != nil {
			s.PipelineName = video.VideoPipeline.Name()
		}
	}

	json.NewEncoder(w).Encode(s)
}

// getCameraStatus will return the real-time information from the camera
func getPipelinesHndl(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(pipelines)
}
