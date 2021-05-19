package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/redeyelab/redeye"
)

type Configuration struct {
	*redeye.Configuration
}

var (
	config Configuration

	cameraList []string
	web        *redeye.WebServer
)

func init() {
	config.Configuration = &redeye.Config
	flag.StringVar(&config.Addr, "addr", ":8000", "Address to serve up redeye from")
	flag.StringVar(&config.Broker, "broker", "tcp://10.24.10.10:1883", "MQTT Broker")
	flag.StringVar(&config.BasePath, "basepath", "/redeye", "BasePath for MQTT Topics and REST URL")
	flag.StringVar(&config.ID, "id", "", "Set the ID for this node")
}

func main() {
	log.Println("Redeye VCR Starting, parsing args...")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("Create and start the new messanger")
	msg := redeye.NewMessanger(config.Broker, config.BasePath)
	msgQ, err := msg.Start()
	if err != nil {
		log.Fatal("Error unable to start messanger, shutting down")
	}

	log.Println("Subscribe to cameras announce ments")
	msg.SubscribeCameras()

	log.Println("Startup our web server ")
	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	// Announce our presence on the camera channel
	log.Println("Announce our Presense")
	msg.Publish(config.BasePath + "/announce/controller/"+msg.Name, msg.Name)

	log.Println("Running the main event loop")
	for true {

		var cmd redeye.TLV
		select {
		case cmd = <-msgQ:
			log.Println("MSG: ", cmd)

		default:
			// log.Println("Main Event Loop, nothing much to do but pause for a moment ...")
			time.Sleep(time.Second * 10)
		}
	}

	wg.Wait()
}
