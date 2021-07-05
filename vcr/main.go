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
	flag.StringVar(&config.BasePath, "basepath", "redeye", "BasePath for MQTT Topics and REST URL")
	flag.StringVar(&config.ID, "id", "", "Set the ID for this node")
}

func main() {

	log.Println("Redeye VCR Starting, parsing args...")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	log.Println("Create and start the new messanger")
	msg := redeye.GetMessanger()
	msgQ, err := msg.Start(&wg)
	if err != nil {
		log.Fatal("Error unable to start messanger, shutting down")
	}

	log.Println("Startup i web server ")
	wg.Add(1)
	web = redeye.GetWebServer(config.Addr, "/" + config.BasePath)
	go web.Start(&wg)

	// Announce our presence on the camera channel
	// topic := "/announce/controller"
	// log.Println("Announce our Presense: ", topic)
	// msg.Publish(topic, msg.Name)

	log.Println("Running the main event loop")
	log.Printf("Subscribers: %+v\n", msg.Subscriptions)
	for true {

		var cmd redeye.TLV
		select {
		case cmd = <-msgQ:
			log.Println("MSG: ", cmd)

		default:
			time.Sleep(time.Millisecond * 10)
		}
	}

	wg.Wait()
}
