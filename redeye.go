package main

import (
	"flag"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	websock    *Websock
	controller string
	webQ       chan interface{}
)

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	SetLogLevel(config.Loglevel)
	startupInfo()

	var wg sync.WaitGroup
	wg.Add(3)

	// Ensure messanger has started, then video play
	go StartMessanger(&wg, &config)
	go StartHTTP(&wg, &config)
	go StartVideo(&wg, &config)

	wg.Wait()
	log.Info().Msg("Good Bye.")
}
