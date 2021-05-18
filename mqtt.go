package redeye

import (
	"log"
	"time"
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
	return messanger
}

// NewMessanger creates a new mqtt messanger
func NewMessanger(broker, path string) (m *Messanger) {
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
func (m *Messanger) Start() (q chan TLV) {

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
	log.Printf("MESSANGER: %+v\n", m)

	// Have the client connect to the broker
	log.Println("Connect to MQTT broker: ", m.Broker)
	if t := m.Client.Connect(); t.Wait() && t.Error() != nil {
		m.Error = t.Error()
		log.Println("error", m.Error.Error())
		return
	}

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
	log.Println("topic: ", topic, " suscribed to topic")
	m.Subscriptions = append(m.Subscriptions, topic)
}

func (m *Messanger) SubscribeCameras() {
	m.Subscribe(m.BasePath + "/announce/camera")
}

// handle all incoming MQTT messages here.
func (m *Messanger) handleIncoming(client mqtt.Client, msg mqtt.Message) {

	topic := msg.Topic()
	payload := string(msg.Payload())

	log.Println("MQTT [In] ", topic, "payload", payload)

	// XXX - This needs to be handled mo betta.
	switch topic {

	case m.BasePath + "/announce/camera":

		log.Println("Creating a new camera! ", payload)		
		cam := NewCamera(payload)
		log.Printf("New camera: %+v\n", cam)


	default:
		log.Println("ERROR - topic", topic, "unknown command")
	}
}

func (m *Messanger) Publish(topic, text string) {

	tstr := m.BasePath + topic 
	if m.Client == nil {
		log.Println("function", "Announce", "Expected client to be connected")
	}

	log.Println("Topic", tstr, "Name", m.Name, "announcing our presence")
	token := m.Client.Publish(tstr, 0, false, m.Name)
	token.Wait()
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

