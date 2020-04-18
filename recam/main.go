package main

import (
	"flag"
	"os"

	"github.com/redeyelab/redeye"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	config *redeye.Settings
	msg    Service
	vid    Service
	web    Service

	cmdQ chan TLV
	msgQ chan TLV
	vidQ chan TLV
	webQ chan TLV
)

func init() {
	cmdQ = make(chan string)
}

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	startupInfo()

	cmdQ = make(chan TLV)

	web = NewWebServer(&config)
	webQ = web.Start(cmdQ)

	msg = NewMessanger(&config)
	msgQ = msg.Start(cmdQ)

	vid = NewVideoPlayer(&config)
	vidQ = vid.Start(cmdQ)

	var cmd TLV
	var src string

	// Accept incoming messages from all running services.
	for cmd != TLVTerm {
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
		switch cmd.Type() {
		case TLVTerm:
			// allow it to exit the outter loop upon the next iteration

		case TLVPlay, TLVPause:
			log.Info().
				Str("dst", "video").
				Str("cmd", cmd.Type()).
				Msg("forwarding message")
			vidQ <- cmd

		default:
			log.Warn().Str("cmd", cmd).Msg("Uknown command...")
		}
	}
	log.Info().Msg("Good Bye.")
}
