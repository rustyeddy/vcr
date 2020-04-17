package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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

var (
	webQ chan string
)

// ===================== Websocket =====================================
func wsUpgradeHndl(w http.ResponseWriter, r *http.Request) {

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Websocket Upgrade failed")
		return
	}
	defer conn.Close()

	log.Info().Msg("ws Upgrader starting reader")
	go wsReader(conn)

	log.Info().Msg("ws Upgrader starting writer")
	wsWriter(conn)

	log.Info().Msg("ws Upgrader leaving")
}

// =============== Websocket Reader =====================================
func wsReader(conn *websocket.Conn) {
	for {
		var msg Message = Message{}
		var err error

		if err = conn.ReadJSON(&msg); err != nil {
			log.Error().Str("error", err.Error()).Msg("reading json msg")
			return
		}

		log.Info().
			Str("type", msg.T).
			Int("len", msg.L).
			Str("Val", msg.V).
			Msg("ws incoming")

		switch msg.T {
		case "ai":
			// if msg.V == "on" {
			// 	video.VideoPipeline, err = GetPipeline(config.Pipeline)
			// 	if err != nil {
			// 		log.Error().
			// 			Str("error", err.Error()).
			// 			Str("pipeline", config.Pipeline).
			// 			Msg("failed to get pipeline")
			// 	}
			// } else if msg.V == "off" {
			// 	video.VideoPipeline = nil
			// }

		case "video":
			// if msg.V == "on" || msg.V == "start" {
			// 	go video.StartVideo()
			// } else if msg.V == "off" || msg.V == "stop" {
			// 	go video.StopVideo()
			// }
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

			if buf, err = json.Marshal(&msg); buf == nil {
				log.Error().Msg("WS unmarshal JSON failed")
				return
			}

			if err != nil {
				log.Error().Msg("WS unmarshal JSON failed")
				return
			}
			log.Debug().Msg("Message sent ... ")
		}
	}
}
