package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Configuration struct {
	Addr     string `json:"addr"`
	Broker   string `json:"broker"`
	Pipeline string `json:"pipeline"`
	Thumb    string `json:"thumb"`
	Vidsrc   string `json:"vidsrc"`
	Vidaddr  string `json:"vidaddr"`
}

var (
	config Configuration
)

func init() {
	flag.StringVar(&config.Addr, "addr", ":8000", "Address to serve up redeye from")
	flag.StringVar(&config.Broker, "broker", "tcp://10.24.10.10:1833", "MQTT Broker")
	flag.StringVar(&config.Thumb, "thumb", "img/thumbnail.jpg", "Thumbnail Image")
	flag.StringVar(&config.Vidsrc, "vidsrc", "0", "Video Source")
	flag.StringVar(&config.Vidsrc, "vidaddr", "8877", "Video Address")
}

func (c *Configuration) Save(path string) (err error) {

	buf, err := json.Marshal(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal JSON from configuration")
		return err
	}

	err = ioutil.WriteFile(path, buf, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed to write config JSON to file")
		return err
	}
	return err
}

// ServeHTTP provides the Web service for the configuration module
func (c *Configuration) ServeHTTP(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(config)
}
