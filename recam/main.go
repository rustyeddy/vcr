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

	cmdQ chan string
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

	var msgQ chan string
	if cfg["msg"] {
		msg = NewMessanger(&config)
		msgQ = msg.Start()
	}

	if cfg["vid"] {
		vid = NewVideoPlayer(&config)
		vid.Start()
	}

	cmdQ := make(chan string)
	var cmd string
	for cmd != "exit" {
		select {
		case cmd = <-webQ:
			log.Info().Str("cmd", cmd).Msg("webQ command")

		case cmd = <-msgQ:
			log.Info().Str("cmd", cmd).Msg("webQ command")

		case cmd = <-cmdQ:
			log.Info().Str("cmd", cmd).Msg("cmdQ command")

		case cmd = <-msgQ:
			log.Info().Str("cmd", cmd).Msg("msgQ command")
		}
	}
	log.Info().Msg("Good Bye.")
}
