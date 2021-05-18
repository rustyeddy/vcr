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
	log.Println("Redeye VCR Starting, parsing args...")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	msg := redeye.NewMessanger(config.Broker, config.BasePath)
	msgQ := msg.Start()
	// msg.SubscribeCameras()
	// msg.Subscribe("/foo")

	web = redeye.NewWebServer(config.Addr, config.BasePath)
	go web.Start(&wg)

	vid = redeye.NewVideoPlayer()
	vidQ = vid.Start(cmdQ)

	vidQ <- redeye.NewTLV(redeye.CMDPlay, 2)

	// // Announce our presence on the camera channel
	msg.Publish("/announce/camera/"+msg.Name, msg.Name)
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
