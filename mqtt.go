package main

import (
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Messanger handles messages and video channels
type Messanger struct {
	Broker         string // MQTT Broker
	ControlChannel string

	mqtt.Client
}

// NewMessager create a New messanger
func NewMessanger(config *Configuration) (msg *Messanger) {
	msg = &Messanger{
		Broker:         config.MQTT,
		ControlChannel: video.GetControlChannel(),
	}
	return msg
}

// Start creates the MQTT client and turns the messanger on
func (m *Messanger) Start(done <-chan interface{}, wg *sync.WaitGroup) {

	opts := mqtt.NewClientOptions().AddBroker(config.MQTT).SetClientID(config.Name)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(m.handleIncoming)
	opts.SetPingTimeout(10 * time.Second)

	m.Client = mqtt.NewClient(opts)
	if t := m.Client.Connect(); t.Wait() && t.Error() != nil {
		err := t.Error()
		log.Fatal(err.Error())
	}

	l.WithFields(log.Fields{
		"broker":  config.MQTT,
		"channel": m.ControlChannel,
	}).Info("Start MQTT Listener")

	if t := m.Client.Subscribe(m.ControlChannel, 0, nil); t.Wait() && t.Error() != nil {
		log.Fatal(t.Error().Error())
	}
	l.WithField("topic", m.ControlChannel).Info("suscribed to topic")
	l.WithField("announce", video.Addr).Info("Announcing Ourselves")
	m.Announce()

	<-done
}

func (m *Messanger) handleIncoming(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	l.WithFields(log.Fields{
		"topic":   topic,
		"message": payload,
	}).Info("MQTT incoming message.")

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
			if video.VideoPipeline == nil {
				video.VideoPipeline = GetPipeline("face")
			} else {
				// Do we need to stop something .?.
				video.VideoPipeline = nil
			}
			// toggle ai
			// video.
			break

		case "hello":
			messanger.Announce()
			break

		default:
			l.WithField("topic", topic).Error("unknown command")
		}
	}
}

// Announce ourselves to the announce channel
func (m *Messanger) Announce() {
	data := video.GetAnnouncement()
	if m.Client == nil {
		log.WithField("function", "Announce").Error("Expected client to be connected")
	}

	log.WithFields(log.Fields{
		"Topic": "camera/announce",
		"Data":  data,
	}).Info("announcing our presence")
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
