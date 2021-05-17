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
	msg.Start()

	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	for (true) {
		time.Sleep(time.Second * 10)

		// Announce our presence on the camera channel
		msg.Publish("/announce/controller/" + msg.Name, msg.Name)
	}

	wg.Wait()
}
