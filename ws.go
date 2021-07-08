package redeye

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSServer struct {
	c websocket.Conn
}

type KeyVal struct {
	K string
	V interface{}
}

type CamerasMsg struct {
	Cameras []*Camera			`json:"cameras"`
	Action string `json:"action"`
}

// ServeHTTP
func (ws WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("Warning Cors Header to '*'")
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:       []string{"echo"},
		InsecureSkipVerify: true, // Take care of CORS
		// OriginPatterns: ["*"],
	})

	if err != nil {
		log.Println("ERROR ", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "houston, we have a problem")

	log.Println("Wait a minute...")
	tQ := time.Tick(time.Second)

	go func() {
		running := true
		for running {

			select {
			case now := <-tQ:
				t := NewTimeMsg(now)
				log.Printf("Sending time %q", t)

				err = wsjson.Write(r.Context(), c, t)
				if err != nil {
					log.Println("ERROR: ", err)
					running = false
				}

				msg := CamerasMsg{ GetCameraList(), "setCameras" }
				err = wsjson.Write(r.Context(), c, msg)
				if err != nil {
					log.Println("ERROR writing cameras: ", err)
					running = false
				}

			}
		}
	}()

	for {
		data := make([]byte, 8192)
		_, data, err := c.Read(r.Context())
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Println("ws Closed")
			return
		}
		if err != nil {
			log.Println("ERROR: reading websocket ", err)
			return
		}
		log.Printf("incoming: %s", string(data))
	}

}

func echo(ctx context.Context, c *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
