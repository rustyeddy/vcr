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
	web    *redeye.WebServer
	vid    *redeye.VideoPlayer

	cmdQ chan redeye.TLV
	vidQ chan redeye.TLV
)

func init() {
	config.Configuration = &redeye.Config
	flag.StringVar(&config.Addr, "addr", ":8000", "Address to serve up redeye from")
	flag.StringVar(&config.Broker, "broker", "tcp://10.24.10.10:1883", "MQTT Broker")
	flag.StringVar(&config.BasePath, "basepath", "/redeye", "BasePath for MQTT Topics and REST URL")
	flag.StringVar(&config.ID, "id", "", "Set the ID for this node")
}

func main() {
	log.Println("Redeye Camera Starting Starting, parsing args...")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("Connect to our message broker")
	msg := redeye.GetMessanger()
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

	log.Println("Subscribe to the Controllers")
	wg.Add(1)
	go msg.SubscribeControllers(&wg)

	log.Println("Announce our Presense")
	msg.Publish("/announce/camera/"+msg.Name, msg.Name)
	log.Printf("Subscribers: %+v\n", msg.Subscriptions)
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
