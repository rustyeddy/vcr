package main

import (
	"flag"
	"log"
	"sync"

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
	wg.Add(2)

	messanger := redeye.NewMessanger(config.Broker, config.BasePath)
	messanger.Start()

	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	wg.Wait()
}
