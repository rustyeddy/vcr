package main

import (
	"flag"
	"os"

	"redeye"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	config *redeye.Settings
	msg    *redeye.Messanger
	vid    *redeye.VideoPlayer
	web    *redeye.WebServer

	cmdQ chan redeye.TLV
	msgQ chan redeye.TLV
	vidQ chan redeye.TLV
	webQ chan redeye.TLV
)

func init() {
	cmdQ = make(chan redeye.TLV)
	d := map[string]string{
		"addr":       ":8000",
		"broker":     "tcp://10.24.10.10:1883",
		"thumb":      "img/thumbnail.jpg",
		"vidsrc":     "0",
		"video-addr": ":8887",
	}
	config = redeye.NewSettings(d)
}

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	//redeye.startupInfo()

	web = redeye.NewWebServer(config)
	webQ = web.Start(cmdQ)

	msg = redeye.NewMessanger(config)
	msgQ = msg.Start(cmdQ)

	vid = redeye.NewVideoPlayer(config)
	vidQ = vid.Start(cmdQ)

	if len(os.Args) > 1 {
		vid.Camstr = os.Args[1]
	}

	var src string
	var cmd redeye.TLV;

	// Accept incoming messages from all running services.
	for cmd.Len() == 0 || cmd.Type() != redeye.CMDTerm {
		log.Info().Msg("Command Q listening for command c.... ")
		select {
		case cmd = <-webQ:
			src = "webQ"
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
		case redeye.CMDTerm:
			// allow it to exit the outter loop upon the next iteration

		case redeye.CMDPlay, redeye.CMDPause:
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
