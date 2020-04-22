package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	config *Settings
	msg    *Messanger
	vid    *VideoPlayer
	web    *WebServer
	mjpg   *MJPEGServer

	cmdQ chan TLV
	msgQ chan TLV
	vidQ chan TLV
	webQ chan TLV

	mjpgQ chan []byte // video frames
)

func init() {
	cmdQ = make(chan TLV)
	d := map[string]string{
		"addr":       ":8000",
		"broker":     "tcp://10.24.10.10:1883",
		"thumb":      "img/thumbnail.jpg",
		"vidsrc":     "0",
		"video-addr": ":8887",
	}
	config = NewSettings(d)
}

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	startupInfo()
	vid = NewVideoPlayer(config)
	vidQ = vid.Start(cmdQ)

	web = NewWebServer(config)
	webQ = web.Start(cmdQ)

	msg = NewMessanger(config)
	msgQ = msg.Start(cmdQ)

	mjpg = NewMJPEGServer(config)
	mjpgQ = mjpg.Start(cmdQ)

	// we have our video camera object set the camstr for this
	// object, we will add it now.
	if len(os.Args) > 1 {
		vid.Camstr = os.Args[1]
	}

	var src string
	cmd := TLV{make([]byte, 2)}

	// Accept incoming messages from all running services.
	for cmd.Type() != TLVTerm {
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
		case TLVTerm:
			// allow it to exit the outter loop upon the next iteration

		case TLVPlay, TLVPause:
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
