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
	cmdQ chan TLV
)

func init() {
	// cmdQ = make(chan TLV)
}

func main() {
	log.Println("Redeye VCR Starting...")

	flag.Parse()

	messanger = NewMessanger()
	messanger.Start(cmdQ)

	var wg sync.WaitGroup
	wg.Add(1)
	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	// var p aeye.Pipeline
	// log.Printf("pipe: %+v\n", p)

	log.Println("Waiting for web to end")
	wg.Wait()
}
