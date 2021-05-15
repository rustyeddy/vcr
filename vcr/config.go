package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

type Configuration struct {
	Addr     string `json:"addr"`
	Broker   string `json:"broker"`
	Pipeline string `json:"pipeline"`
	Thumb    string `json:"thumb"`
	Vidsrc   string `json:"vidsrc"`
	Vidaddr  string `json:"vidaddr"`
}

func init() {
	flag.StringVar(&config.Addr, "addr", ":8000", "Address to serve up redeye from")
	flag.StringVar(&config.Broker, "broker", "tcp://10.24.10.10:1833", "MQTT Broker")
}

func (c *Configuration) Save(path string) (err error) {

	buf, err := json.Marshal(c)
	if err != nil {
		log.Println("failed to marshal JSON from configuration")
		return err
	}

	err = ioutil.WriteFile(path, buf, 0644)
	if err != nil {
		log.Println("failed to write config JSON to file")
		return err
	}
	return err
}

// ServeHTTP provides the Web service for the configuration module
func (c *Configuration) ServeHTTP(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(config)
}
