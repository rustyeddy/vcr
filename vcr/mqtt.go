package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Messanger handles messages and video channels
type Messanger struct {
	Name          string
	Broker        string
	Subscriptions []string

	mqtt.Client
	Error error
}

var (
	messanger *Messanger
)

func GetMessanger() *Messanger {
	if messanger == nil {
		messanger = &Messanger{
			Name:          "mqtt",
			Broker:        config.Broker,
			Subscriptions: []string{"camera/announce"},
		}
	}
	return messanger
}

// NewMessanger creates a new mqtt messanger
func NewMessanger() (m *Messanger) {
	m = &Messanger{
		Name:   GetHostname(),
		Broker: config.Broker,
	}

	if m.Name == "" {
		log.Fatal("Expected a hostname got (nil)")
	}
	sub := "camera/" + m.Name
	m.Subscriptions = []string{sub}
	return m
}

// Start fires up our MQTT client, then subscribes the given subscription
// list.
func (m *Messanger) Start(cmdQ chan TLV) (q chan TLV) {

	// set up the MQTT client options
	opts := mqtt.NewClientOptions().AddBroker(m.Broker).SetClientID(m.Name)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(m.handleIncoming)
	opts.SetPingTimeout(10 * time.Second)

	// create a NewClient
	m.Client = mqtt.NewClient(opts)
	if m.Client == nil {
		log.Println("New Client Failed, no MQTT available")
		return
	}

	// Have the client connect to the broker
	log.Println("Connect to MQTT broker: ", m.Broker)
	if t := m.Client.Connect(); t.Wait() && t.Error() != nil {
		m.Error = t.Error()
		log.Println("error", m.Error.Error())
		return
	}

	// Roll through the subscription list and subscribe
	for _, topic := range m.Subscriptions {
		log.Println("topic", topic, " MQTT Subscribe to topic...")
		m.Subscribe(topic)
	}

	// Announce our presence on the camera channel
	m.Announce()

	q = make(chan TLV)
	log.Println("messanger gofuncing listener")
	go func() {
		for {
			log.Println("Waiting for message ... ")
			select {
			case cmd := <-q:
				log.Println("cmd", cmd.Str(), "got a message.")
				switch cmd.Type() {
				case CMDTerm:
					log.Println("Exiting messanger")
					return
				default:
					log.Println("cmd is not supported")
				}
			}
		}
	}()
	return q
}

// Subscribe to the given channel
func (m *Messanger) Subscribe(topic string) {
	log.Println("broker ", m.Broker, " channel ", topic, " Start MQTT Listener")

	if t := m.Client.Subscribe(topic, 0, nil); t.Wait() && t.Error() != nil {
		log.Println("error", t.Error().Error(), "Failed to subscribe to mqtt socket")
		return
	}
	log.Println("topic", topic, " suscribed to topic")
	m.Subscriptions = append(m.Subscriptions, topic)
}

// handle all incoming MQTT messages here.
func (m *Messanger) handleIncoming(client mqtt.Client, msg mqtt.Message) {

	topic := msg.Topic()
	payload := string(msg.Payload())

	log.Println("topic", topic, "payload", payload, "MQTT incoming message.")

	switch {
	case strings.Compare(topic, "camera/announce") == 0:

		m.Announce()

	case strings.Contains(topic, "camera/"):
		switch payload {

		case "pause", "off", "play", "on":
			log.Println("msg", payload, "sending payload to cmdQ")
			buf := make([]byte, 2)
			buf[1] = 2 // all our messages are two bytes!
			switch payload {
			case "on", "play":
				buf[0] = CMDPlay
			case "off", "pause":
				buf[0] = CMDPause
			case "snap":
				buf[0] = CMDSnap
			default:
				log.Println("str", payload, "Unsupported Msg Type")
				return
			}

			// cmdQ <- TLV{buf}
			break

		case "ai":
			// TODO Send a message rather than handle here

			/*
				 var err error

				if video.VideoPipeline == nil {
					video.VideoPipeline, err = GetPipeline(config.Pipeline)
					if err != nil {
						log.Println().Str("pipeline", config.Pipeline).Msg("Failed to get pipeline")
						return
					}
				} else {
					// Do we need to stop something .?.
					video.VideoPipeline = nil
				}
			*/
			break

		case "hello":

			// Announce ourselves
			m.Announce()
			break

		default:
			log.Println("topic", topic, "unknown command")
		}
	}
}

// TODO Move this to video ...
// Announce ourselves to the announce channel
func (m *Messanger) Announce() {

	//data := "camera/" + GetHostname()
	topic := config.BasePath + "/announce/controller/" + m.Name
	if m.Client == nil {
		log.Println("function", "Announce", "Expected client to be connected")
	}

	log.Println("Topic", topic, "Name", m.Name, "announcing our presence")
	token := m.Client.Publish(topic, 0, false, m.Name)
	token.Wait()
}

// getMessanger
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var status *MessangerStatus
	if m := GetMessanger(); m != nil {
		status = m.GetStatus()
	} else {
		// serve up the null entry
		status = &MessangerStatus{
			Broker: "DISCONNECTED",
		}
	}
	json.NewEncoder(w).Encode(status)
}

// MessangerSstatus returns the status of the currently
// running Messanger.
type MessangerStatus struct {
	Broker        string
	Subscriptions []string
	Connected     bool
}

// GetMessangerStatus lets the caller know what is happening
// with the messanger.
func (m *Messanger) GetStatus() (ms *MessangerStatus) {
	ms = &MessangerStatus{
		Broker:        m.Broker,
		Subscriptions: m.Subscriptions,
	}
	ms.Connected = m.Client != nil
	return ms
}
