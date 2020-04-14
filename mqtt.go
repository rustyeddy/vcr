package main

import (
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

var (
	messanger *Messanger
)

// Messanger handles messages and video channels
type Messanger struct {
	Broker         string
	Subscriptions []string

	mqtt.Client
	Error error
}

func GetMessanger() *Messanger {
	if messanger == nil {
		messanger = &Messanger{
			Broker: config.MQTT,
			Subscriptions: nil,
		}
	}
	return messanger
}

// StartMessanger 
func StartMessanger(wg *sync.WaitGroup, config *Configuration) {
	defer wg.Done()

	m := GetMessanger()
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

	log.Info().Str("announce", video.Addr).Msg("Announcing Ourselves")
	m.Announce()
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
		controller = payload
		m.Announce()

	case strings.Contains(topic, "camera/"):
		switch payload {

		case "on":
			go video.StartVideo()
			break

		case "off":
			video.StopVideo()
			break

		case "ai":
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
			break

		case "hello":
			m.Announce()
			break

		default:
			log.Error().Str("topic", topic).Msg("unknown command")
		}
	}
}

// Announce ourselves to the announce channel
func (m *Messanger) Announce() {
	data := video.GetAnnouncement()
	if m.Client == nil {
		log.Error().Str("function", "Announce").Msg("Expected client to be connected")
	}

	log.Info().
		Str("Topic", "camera/announce").
		Str("Data", data).
		Msg("announcing our presence")
	token := m.Client.Publish("camera/announce", 0, false, data)
	token.Wait()
}

// Read stuff
func (m *Messanger) Read(b []byte) (n int, err error) {
	panic("Implement reader")
	return n, err
}

// Write stuff
func (m *Messanger) Write(b []byte) (n int, err error) {
	panic("Implement writer")
	return n, err
}
