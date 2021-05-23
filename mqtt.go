package redeye

import (
	"fmt"
	"log"
	"time"
	"sync"
	"os"

	"encoding/json"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Messanger handles messages and video channels
type Messanger struct {
	Name          string
	Broker        string
	BasePath	  string

	Subscriptions []string
	Published	  []string

	mqtt.Client
	Error error
}

var (
 	messanger *Messanger
)

func GetMessanger() *Messanger {
	if messanger == nil {
		log.Println("Creating New messanger ", Config.Broker, Config.BasePath)
		messanger = newMessanger(Config.Broker, Config.BasePath)
	}
 	return messanger
}

// NewMessanger creates a new mqtt messanger
func newMessanger(broker, path string) (m *Messanger) {
	m = &Messanger{
		Name:   GetHostname(),
		Broker: broker,
		BasePath: path,
	}

	if m.Name == "" {
		log.Fatal("Expected a hostname got (nil)")
	}
	return m
}

// Start fires up our MQTT client
func (m *Messanger) Start(wg *sync.WaitGroup) (q chan TLV, err error) {

	mqtt.WARN = log.New(os.Stdout, "[DEBUG] ", 0)
	
	// set up the MQTT client options
	opts := mqtt.NewClientOptions().AddBroker(m.Broker).SetClientID(m.Name)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(m.handleIncoming)
	opts.SetPingTimeout(10 * time.Second)
	
	// create a NewClient
	m.Client = mqtt.NewClient(opts)
	if m.Client == nil {
		return nil, fmt.Errorf("New Client Failed, no MQTT available")
	}
	// Have the client connect to the broker
	if t := m.Client.Connect(); t.Wait() && t.Error() != nil {
		m.Error = t.Error()
		return nil, fmt.Errorf("Error connecting to MQTT broker: %w", m.Error)
	}

	log.Println("Args: ", os.Args[0])

	if (os.Args[0] == "./vcr") {
		m.SubscribeCameras()		
	} else {
		m.SubscribeControllers()		
	}

	q = make(chan TLV)
	go func() {
		for {
			log.Println("Waiting for an incoming MQTT cmd")
			select {
			case cmd := <-q:
				log.Println("cmd", cmd.Str(), "got a message.")
				switch cmd.Type() {
				case CMDTerm:
					return
				default:
					if Config.Debug {
						log.Println("cmd is not supported")						
					}
				}
			}
		}
	}()
	return q, nil
}

// Subscribe to the given channel
func (m *Messanger) Subscribe(topic string) error {

	log.Println("Subscribe broker ", m.Broker, "topic: ", topic)
	if t := m.Client.Subscribe(topic, 0, nil); t.Wait() && t.Error() != nil {
		return fmt.Errorf("Failed to subscribe to mqtt socket: %w", t.Error())
	} else {
		log.Printf("Subscription succeeded %s", topic)
	}
	m.Subscriptions = append(m.Subscriptions, topic)
	return nil
}


func (m *Messanger) SubscribeCameras() error {
	topic := m.BasePath + "/announce/camera"
	return m.Subscribe(topic)
}

func (m *Messanger) SubscribeControllers() error {
	topic := m.BasePath + "/announce/controller"
	return m.Subscribe(topic)
}

// handle all incoming MQTT messages here.
func (m *Messanger) handleIncoming(client mqtt.Client, msg mqtt.Message) {

	topic := msg.Topic()
	payload := string(msg.Payload())

	log.Println("MQTT [In] ", topic, "payload", payload)		
	switch topic {

	case m.BasePath + "/announce/controller":
		cam := NewCamera(payload)
		fmt.Printf("Controller: %+v\n", cam)

	case m.BasePath + "/announce/camera":
		cam := NewCamera(payload)
		fmt.Printf("Camera: %+v\n", cam)

	default:
		log.Println("Incoming Message - topic ", topic, ", payload: ", payload)
	}
}

func (m *Messanger) Publish(topic, text string) error {

	tstr := m.BasePath + topic 
	if m.Client == nil {
		return fmt.Errorf("Failed to Publish topic %s", topic)
	}

	log.Println("Publishing topic: ", tstr, " payload ", m.Name)
	token := m.Client.Publish(tstr, 0, false, m.Name)
	token.Wait()
	log.Println("Publishing topic: ", tstr, " payload ", m.Name, " done waiting")
	return nil
}

// getMessanger
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var status *MessangerStatus
	if m := messanger; m != nil {
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

