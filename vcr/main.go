package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/redeyelab/redeye"
	//"github.com/redeyelab/redeye/aeye"
)

var (
	config Configuration

	cameraList []string
	web *redeye.WebServer;
)

func main() {
	log.Println("Redeye VCR Starting, parsing args...")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	msg := redeye.NewMessanger(config.Broker, config.BasePath)
	msgQ := msg.Start()
	msg.SubscribeCameras()
	msg.Subscribe("/foo")

	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	// Announce our presence on the camera channel
	msg.Publish("/announce/controller/" + msg.Name, msg.Name)

	for (true) {

		var cmd redeye.TLV
		select {
		case cmd = <-msgQ:
			log.Println("MSG: ", cmd)

		default:
			log.Println("Main Event Loop, nothing much to do but pause for a moment ...")
			time.Sleep(time.Second * 10)
		}
	}

	wg.Wait()
}
