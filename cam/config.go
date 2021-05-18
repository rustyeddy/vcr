package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

// Configuration struct handles the startup and running configuration for
// redeye vcr, includes reading and writing to file.
type Configuration struct {
	ID		 string `json:"id"`		  // The ID
	Addr     string `json:"addr"`	  // the address we'll serve up 
	Broker   string `json:"broker"`	  // MQTT broker address
	BasePath string `json:"basepath"` // /redeye
}

func init() {
	flag.StringVar(&config.Addr, "addr", ":8000", "Address to serve up redeye from")
	flag.StringVar(&config.Broker, "broker", "tcp://localhost:1883", "MQTT Broker")
	flag.StringVar(&config.BasePath, "basepath", "/redeye", "BasePath for MQTT Topics and REST URL")
	flag.StringVar(&config.ID, "id", "", "Set the ID for this node")
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
