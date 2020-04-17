package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	cfg map[string]bool
	msg *Messanger
	vid *VideoPlayer
	web *WebServer

	cmdQ chan string // incoming commands to cmd xchange
	msgQ chan string // outgoing mqtt
	vidQ chan string // control video stream
)

func init() {

	// TODO make these flags
	cfg = map[string]bool{
		"web": true,
		"msg": true,
		"vid": true,
	}

	cmdQ = make(chan string)
}

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	startupInfo()

	// Create and configure all the services
	if cfg["web"] {
		web = NewWebServer(&config)
		web.Start()
	}

	msgQ = make(chan string)
	if cfg["msg"] {
		msg = NewMessanger(&config)
		msgQ = msg.Start()
	}

	vidQ = make(chan string)
	if cfg["vid"] {
		vid = NewVideoPlayer(&config)
		vidQ = vid.Start()
	}

	cmdQ = make(chan string)
	var cmd string
	var src string

	// Accept incoming messages from all running services.
	for cmd != "exit" {
		log.Info().Msg("Command Q listening for command c.... ")
		select {
		case cmd = <-webQ:
			src = "webQ"

		case cmd = <-msgQ:
			src = "msgQ"

		case cmd = <-cmdQ:
			src = "cmdQ"
		}

		log.Info().
			Str("src", src).
			Str("cmd", cmd).
			Msg("Command Exchange Incoming")

		// Send the command off to any reciever
		switch cmd {
		case "exit":
			// allow it to exit the outter loop upon the next iteration

		case "play", "on", "pause", "off":
			log.Info().
				Str("dst", "video").
				Str("cmd", cmd).
				Msg("forwarding message")
			vidQ <- cmd

		default:
			log.Warn().Str("cmd", cmd).Msg("Uknown command...")
		}
	}
	log.Info().Msg("Good Bye.")
}
