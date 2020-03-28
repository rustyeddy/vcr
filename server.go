package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/hybridgroup/mjpeg"
)

type Server struct {
	// HTTP Server and router handle all HTML and Static page requests
	// including all CSS and JavaScript files.
	*http.Server
	*mux.Router
}

// NewServer will return a new server.!.
func NewServer(config *Configuration) (srv *Server) {
	srv = &Server{}
	return srv
}

// StartServer is a meta starter, it starts the HTTP for the SPA and
// the MJPEG server. TODO: should /mjpeg and index.html be on the same
// server (port)?
func (srv *Server) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	srv.Router = mux.NewRouter()
	srv.Server = &http.Server{
		Handler: srv.Router,
		Addr:    config.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// ----------------------------------------------
	srv.Router.HandleFunc("/ws", wsUpgradeHndl)
	srv.Router.HandleFunc("/api/health", healthCheckHndl)
	srv.Router.HandleFunc("/api/config", getConfigHndl)
	srv.Router.HandleFunc("/api/config/{key}/{val}", setConfigValsHndl)

	// app = spaHandler{StaticPath: "pub", IndexPath: "index.html"}
	srv.Router.PathPrefix("/").Handler(srv)

	// Set the route for video
	vpath := "/mjpeg"
	l.WithFields(log.Fields{
		"address": config.VideoAddr,
		"path":    vpath,
	}).Info("Start Video Server")
	video.Stream = mjpeg.NewStream()
	http.Handle(vpath, video.Stream)

	// Listen to requests for video at videoaddr. NOTE: VideoServer actually
	// turns the video camera stream on and off through the control api
	// via REST, MQTT or WebUI.
	go http.ListenAndServe(config.VideoAddr, nil)

	// Startup HTTP server
	l.WithField("addr", srv.Addr).Info("Starting HTTP Server ...")
	err := srv.Server.ListenAndServe()
	l.Fatal(err.Error())
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
	log.Fatal("TODO config vals")
}
