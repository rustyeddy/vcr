package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

// ===================== Websocket ===========================
type Message struct {
	T string
	L int
	V string
}

type Websock struct {
}

// ===================== Websocket =====================================
func wsUpgradeHndl(w http.ResponseWriter, r *http.Request) {

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Errorf("Websocket Upgrade failed %v", err)
		return
	}
	defer conn.Close()

	l.Info("ws Upgrader starting reader")
	go wsReader(conn)

	l.Info("ws Upgrader starting writer")
	wsWriter(conn)

	l.Info("ws Upgrader leaving")
}

// =============== Websocket Reader =====================================
func wsReader(conn *websocket.Conn) {
	for {
		var msg Message = Message{}
		var err error

		if err = conn.ReadJSON(&msg); err != nil {
			l.WithError(err)
			return
		}

		l.WithFields(log.Fields{
			"type": msg.T,
			"len":  msg.L,
			"Val":  msg.V,
		}).Infof("ws incoming")
		switch msg.T {
		case "ai":
			if msg.V == "on" {
				video.VideoPipeline, err = GetPipeline(config.Pipeline)
				if err != nil {
					l.WithError(err)
				}
			} else if msg.V == "off" {
				video.VideoPipeline = nil
			}

		case "video":
			if msg.V == "on" || msg.V == "start" {
				go video.StartVideo()
			} else if msg.V == "off" || msg.V == "stop" {
				go video.StopVideo()
			}
		}
	}
}

// =============== Websocket Writer =====================================
func wsWriter(conn *websocket.Conn) {
	for {
		var err error
		select {
		case msg := <-webQ:
			var buf []byte
			log.Debugf("WS Send JSON %+v", msg)

			if buf, err = json.Marshal(&msg); buf == nil {
				log.Error("WS unmarshal JSON failed")
				return
			}

			if err != nil {
				log.Error("WS unmarshal JSON failed")
				return
			}
			log.Debug("Message sent ... ")
		}
	}
}
