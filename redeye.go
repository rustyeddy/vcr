package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	startupInfo()

	// Create and configure all the services
	web := NewWebServer(&config)
	msg := NewMessanger(&config)
	vid := NewVideoPlayer(&config)

	// Start the services
	web.Start()
	vid.Start()
	msgQ := msg.Start()

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

		case cmd = <-cmdQ:
			log.Info().Str("cmd", cmd).Msg("cmdQ command")

		case cmd = <-msgQ:
			log.Info().Str("cmd", cmd).Msg("msgQ command")
		}
	}
	log.Info().Msg("Good Bye.")
}
