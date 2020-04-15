package main

import (
	"flag"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting redeye")

	startupInfo()

	srv := NewHTTPServer(&config)
	//msg := NewMQTTServer(&config)
	//vid := NewVideoServer(&config)

	var wg sync.WaitGroup
	wg.Add(1)
	srv.Start(&wg)
	//msgQ := msg.Start(&wg)
	//vidQ := vid.Start(&wg)

	// Ensure messanger has started, then video play
	//go StartMessanger(&wg, &config)
	//go StartHTTP(&wg, &config)
	//go StartVideo(&wg, &config)

	select {
	case cmd := <-srv.Q:
		log.Info().Str("cmd", cmd).Msg("command")
	}

	wg.Wait()
	log.Info().Msg("Good Bye.")
}
