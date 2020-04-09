package main

import (
	"flag"
	"sync"
)

var (
	config     *Configuration
	server     *Server
	video      *VideoPlayer
	messanger  *Messanger
	websock    *Websock
	controller string

	webQ chan interface{}
)

func init() {
	// Load a config file if we have one
	config = GetConfig()
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Set the logger level
	SetLogLevel(config.Loglevel)

	// Now a done channel to pass to all services
	done := make(chan interface{})

	// start up the HTTP website and REST API. We will start up
	// some more go routine listeners
	var wg sync.WaitGroup

	// Print some start up info the the master desires
	startupInfo()

	// TODO move this elsewhere, we can keep this variable in
	// cases we donot have a video camera. This camera is literally
	// nothing at the moment, it needs to
	video = NewVideoPlayer(config)

	// server provides the static website and REST APIs
	//
	// HTTP static website from pub/index.html (configurable)
	// HTTP REST '/api/...' see respective sections
	// Websocket server also listens for websocket requests
	wg.Add(2)
	server = NewServer(config)
	go server.Start(&wg)

	// messengar will handle realtime messaging
	messanger = NewMessanger(config)
	go messanger.Start(done, &wg)

	// Wait forever or until all of messanger and server fail
	wg.Wait()
	l.Info("Good Bye.")
}
