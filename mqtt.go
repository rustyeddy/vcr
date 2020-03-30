package main

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Messanger struct {
	Broker         string // MQTT Broker
	ControlChannel string

	mqtt.Client
}

func NewMessanger(config *Configuration) (msg *Messanger) {
	camstr := "/camera/"
	if hname, err := os.Hostname(); err != nil {
		l.WithError(err).Fatal("Good bye cruel world!")
	} else {
		camstr += hname
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

		default:
			l.WithField("topic", topic).Error("unknown command")
		}
	}
}

func (m *Messanger) Read(b []byte) (n int, err error) {
	panic("Implement reader")
	return n, err
}

func (m *Messanger) Write(b []byte) (n int, err error) {
	panic("Implement writer")
	return n, err
}
