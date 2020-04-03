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

func NewMessanger(config *Configuration) (msg *Messanger) {
	camstr := "camera/" + GetHostname()
	if len(camstr) <= len("camera/") {
		l.WithField("camera channel", camstr).Error("hostname bad")
	}
	msg = &Messanger{
		Broker:         config.MQTT,
		ControlChannel: camstr,
	}
	return msg
}

func (m *Messanger) Start(done <-chan interface{}, wg *sync.WaitGroup) {

	opts := mqtt.NewClientOptions().AddBroker(config.MQTT).SetClientID(config.Name)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(m.handleIncoming)
	opts.SetPingTimeout(10 * time.Second)

	c := mqtt.NewClient(opts)
	if t := c.Connect(); t.Wait() && t.Error() != nil {
		err := t.Error()
		log.Fatal(err.Error())
	}

	l.WithFields(log.Fields{
		"broker":  config.MQTT,
		"channel": m.ControlChannel,
	}).Info("Start MQTT Listener")

	if t := c.Subscribe(m.ControlChannel, 0, nil); t.Wait() && t.Error() != nil {
		log.Fatal(t.Error().Error())
	}
	l.WithField("topic", m.ControlChannel).Info("suscribed to topic")
	l.WithField("announce", video.Name).Info("Announcing Ourselves")

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
	case strings.Contains(topic, "/camera/"):
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
	ipaddr := GetIPAddr()
	data := ipaddr + ":" + video.Name
	if m.Client != nil {
		log.WithFields(log.Fields{
			"Topic": "camera/announce",
			"Data":  data,
		}).Info("announcing our presence")
		token := m.Client.Publish("camera/announce", 0, false, data)
		token.Wait()
	}
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
