package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	msg *Messanger
	vid *VideoPlayer

	cmdQ chan TLV
	msgQ chan TLV
	vidQ chan TLV
)

func init() {
	cmdQ = make(chan TLV)
}

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	go web()

	msg = NewMessanger()
	msgQ = msg.Start(cmdQ)

	vid = NewVideoPlayer()
	vidQ = vid.Start(cmdQ)

	if len(os.Args) > 1 {
		vid.Camstr = os.Args[1]
	}

	var src string
	var cmd TLV

	vidQ <- NewTLV(CMDPlay, 2)

	// Accept incoming messages from all running services.
	for cmd.Len() == 0 || cmd.Type() != CMDTerm {

		log.Info().Msg("Command Q listening for command c.... ")
		select {
		case cmd = <-msgQ:
			src = "msgQ"

		case cmd = <-cmdQ:
			src = "cmdQ"

		case cmd = <-vidQ:
			src = "vidQ"
		}

		log.Info().
			Str("src", src).
			Str("cmd", cmd.Str()).
			Msg("Command Exchange Incoming")

		// Send the command off to any reciever
		switch cmd.Type() {
		case CMDTerm:
			// allow it to exit the outter loop upon the next iteration

		case CMDPlay, CMDPause:
			log.Info().
				Str("dst", "video").
				Str("cmd", cmd.Str()).
				Msg("forwarding message")
			vidQ <- cmd

		default:
			log.Warn().Str("cmd", cmd.Str()).Msg("Uknown command...")
		}
	}
	log.Info().Msg("Good Bye.")
}
