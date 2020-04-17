package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	serviceConfig map[string]bool
)

func init() {

	// TODO make these flags
	serviceConfig = map[string]bool{
		"web": true,
		"fs":  false,
		"msg": true,
		"vid": true,
	}
}

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	startupInfo()

	// Create and configure all the services
	if serviceConfig["web"] {
		web := NewWebServer(&config)
		web.Start()
	}

	if serviceConfig["fs"] {
		fs := NewFileServer(&config)
		fs.Start()
	}

	var msgQ chan string
	if serviceConfig["msg"] {
		msg := NewMessanger(&config)
		msgQ = msg.Start()
	}

	if serviceConfig["vid"] {
		vid := NewVideoPlayer(&config)
		vid.Start()
	}

	// TODO Have the video player announce itself when msgQ is alive
	//
	// If the messanger is running, subscribe to our control topic
	// if m := GetMessanger(); m != nil {
	// 	m.Subscribe(video.GetControlChannel())
	// }

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
