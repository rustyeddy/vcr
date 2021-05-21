package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/redeyelab/redeye"
)

var (
	config Configuration
	web    *redeye.WebServer
	vid    *redeye.VideoPlayer

	cmdQ chan redeye.TLV
	vidQ chan redeye.TLV
)

func main() {
	log.Println("Redeye Camera Starting Starting, parsing args...")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("Connect to our message broker")
	msg := redeye.NewMessanger(config.Broker, config.BasePath)
	msgQ, err := msg.Start()
	if err != nil {
		log.Fatal("Error connecting to message broker")
	}
	log.Println("Fire up the web server")
	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	log.Println("Grab a new video player and ready it to stream video")
	vid = redeye.NewVideoPlayer()
	vidQ = vid.Start(cmdQ)

	log.Println("Announce our Presense")
	msg.Publish("/announce/camera/"+msg.Name, msg.Name)

	vidQ <- redeye.NewTLV(redeye.CMDPlay, 2)
	for true {

		var cmd redeye.TLV
		select {
		case cmd = <-msgQ:
			log.Printf("msgQ: %+v\n", cmd)

		case cmd = <-vidQ:
			log.Printf("vidQ: %+v\n", cmd)

		default:
			log.Println("Main Event Loop, nothing much to do but pause for a moment ...")
			time.Sleep(time.Second * 10)
		}
	}

	wg.Wait()

}
