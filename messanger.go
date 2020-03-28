package main

import (
	"sync"
	"time"

	"github.com/apex/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Messanger struct {
	Broker string   // MQTT Broker
	Topics []string // List of topics to subscribe to

	mqtt.Client
}

func NewMessanger(config *Configuration) (msg *Messanger) {
	msg = &Messanger{
		Broker: "http://10.24.10.10/",
	}
	return msg
}

func (m *Messanger) Start(done <-chan interface{}, wg *sync.WaitGroup) {

	opts := mqtt.NewClientOptions().AddBroker(config.MQTT).SetClientID("redeye")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(10 * time.Second)

	c := mqtt.NewClient(opts)
	if t := c.Connect(); t.Wait() && t.Error() != nil {
		err := t.Error()
		log.Fatal(err.Error())
	}

	l.WithField("broker", config.MQTT).Info("Start MQTT Listener")

	for _, topic := range m.Topics {
		if t := c.Subscribe(topic, 0, nil); t.Wait() && t.Error() != nil {
			log.Fatal(t.Error().Error())
		}
		l.WithField("topic", topic).Info("suscribed to topic")
	}

	<-done
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	l.WithFields(log.Fields{
		"topic":   topic,
		"message": payload,
	}).Info("MQTT incoming message.")

	switch topic {
	case "/camera/control":
		switch payload {

		case "on":
			go video.StartVideo()
			break

		case "off":
			video.StopVideo()
			break
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

func (m *Messanger) Subscribe(topic string) {

	panic("subscribe to topic")
	m.Topics = append(m.Topics, topic)
}
