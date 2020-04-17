package main

import (
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

var (
	messanger *Messanger
)

// Messanger handles messages and video channels
type Messanger struct {
	Broker        string
	Subscriptions []string

	cmdQ chan string
	mqtt.Client
	Error error
}

func NewMessanger(config *Configuration) *Messanger {
	messanger = &Messanger{
		Broker:        config.MQTT,
		Subscriptions: nil,
	}
	return messanger
}

// StartMessanger
func (m *Messanger) Start() (q chan string) {
	m.cmdQ = q
	opts := mqtt.NewClientOptions().AddBroker(config.MQTT).SetClientID(config.Name)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(m.handleIncoming)
	opts.SetPingTimeout(10 * time.Second)

	m.Client = mqtt.NewClient(opts)
	if m.Client == nil {
		// XXX should we have a connect retry?
		log.Error().Msg("New Client Failed, no MQTT available")
		return
	}

	if t := m.Client.Connect(); t.Wait() && t.Error() != nil {
		m.Error = t.Error()
		log.Error().Str("error", m.Error.Error()).Msg("Failed opening MQTT client")
		return
	}

	for _, topic := range m.Subscriptions {
		m.Subscribe(topic)
	}
	//m.Announce()

	// Start listening for commands to relay on channels
	go func() {
		for {
			cmd := <-q
			switch cmd {
			case "":
				log.Warn().Msg("cmd is empty")
			case "exit":
				log.Info().Msg("Exiting messanger")
				return
			}
		}
	}()
	return
}

// Subscribe to the given channel
func (m *Messanger) Subscribe(topic string) {
	log.Info().
		Str("broker", config.MQTT).
		Str("channel", topic).
		Msg("Start MQTT Listener")

	if t := m.Client.Subscribe(topic, 0, nil); t.Wait() && t.Error() != nil {
		log.Error().Str("error", t.Error().Error()).Msg("Failed to subscribe to mqtt socket")
		return
	}
	log.Info().Str("topic", topic).Msg("suscribed to topic")
	m.Subscriptions = append(m.Subscriptions, topic)
}

func (m *Messanger) handleIncoming(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	log.Info().
		Str("topic", topic).
		Str("message", payload).
		Msg("MQTT incoming message.")

	switch {
	case strings.Compare(topic, "camera/announce") == 0:
		// Ignore the controller
		// m.Announce()

	case strings.Contains(topic, "camera/"):
		switch payload {

		case "on":
			// TODO Create and start using TLVs, but just a string for now
			// go video.StartVideo() -> send message instead
			break

		case "off":
			//video.StopVideo() -> send message instead
			break

		case "ai":
			// TODO Send a message rather than handle here

			/*
				 var err error

				if video.VideoPipeline == nil {
					video.VideoPipeline, err = GetPipeline(config.Pipeline)
					if err != nil {
						log.Error().Str("pipeline", config.Pipeline).Msg("Failed to get pipeline")
						return
					}
				} else {
					// Do we need to stop something .?.
					video.VideoPipeline = nil
				}
			*/
			break

		case "hello":
			//m.Announce()
			break

		default:
			log.Error().Str("topic", topic).Msg("unknown command")
		}
	}
}

// TODO Move this to video ...
// Announce ourselves to the announce channel
// func (m *Messanger) Announce() {
// 	data := video.GetAnnouncement()
// 	if m.Client == nil {
// 		log.Error().Str("function", "Announce").Msg("Expected client to be connected")
// 	}

// 	log.Info().
// 		Str("Topic", "camera/announce").
// 		Str("Data", data).
// 		Msg("announcing our presence")
// 	token := m.Client.Publish("camera/announce", 0, false, data)
// 	token.Wait()
// }
