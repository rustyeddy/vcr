package main

import (
	"flag"
	"log"
	"sync"

	//"github.com/redeyelab/redeye/aeye"
)

var (
	config Configuration
	cmdQ chan TLV
)

func init() {
	// cmdQ = make(chan TLV)
}

func main() {
	log.Println("Redeye VCR Starting...")

	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	go web(wg)

	messanger = NewMessanger()
	messanger.Start(cmdQ)

	// var p aeye.Pipeline
	// log.Printf("pipe: %+v\n", p)

	log.Println("Waiting for web to end")
	wg.Wait()
}
