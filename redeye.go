package main

import (
	"flag"
	"net"
	"os"
	"sync"

	"github.com/apex/log"
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

	//
	// We are going to start the following servers:
	//
	// HTTP static website from pub/index.html (configurable)
	// HTTP REST '/api/...' see respective sections
	// Websocket server also listens for websocket requests
	//
	wg.Add(2)
	server = NewServer(config)
	go server.Start(&wg)

	// MQTT Client connected to /topic/tempf (TODO change channels)
	messanger = NewMessanger(config)
	go messanger.Start(done, &wg)

	// Wait forever or until all of messanger and server fail
	wg.Wait()
	l.Info("Good Bye.")
}

func startupInfo() {
	if !config.Debug {
		return
	}
	log.Infof("config %v\n", config)

	l.WithFields(log.Fields{
		"app":      "redeye",
		"pid":      os.Getpid(),
		"hostname": GetHostname(),
	}).Info("App is starting up ...")
}

// GetHostname for ourselves
func GetHostname() (hname string) {
	var err error
	if hname, err = os.Hostname(); err != nil {
		log.WithError(err).Fatal("Good bye cruel world!")
	}
	return hname
}

// GetIPAddr of ourselves
func GetIPAddr() (addr string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.WithField("addr", addr).Fatal(err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
