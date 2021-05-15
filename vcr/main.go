package main

import (
	"flag"
	"log"
	"sync"
)

var (
	config Configuration
)

func main() {
	log.Println("Redeye VCR Starting...")

	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	go web(wg)

	log.Println("Waiting for web to end")
	wg.Wait()
}
