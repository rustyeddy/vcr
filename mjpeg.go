package main

import (
	"net/http"

	"github.com/hybridgroup/mjpeg"
	"github.com/rs/zerolog/log"
)

// MJPEGServer starts up an HTTP server for posting rapidly
// updating JPEG images from the camera.
type MJPEGServer struct {
	Name          string
	Addr          string
	*mjpeg.Stream `json:"-"` // Stream will always be available
}

// NewMJPEGServer will create a new video player with default nil set.
func NewMJPEGServer(config *Settings) (m *MJPEGServer) {
	m = &MJPEGServer{"mjpg", ":8887", nil}
	return m
}

// Start the mjpeg server clients will likely be either the web client
// or the CV pipeline.
func (m *MJPEGServer) Start(cmdQ chan TLV) (mpgQ chan []byte) {

	// Set the route for video
	mpath := "/mjpeg"
	log.Info().
		Str("address", config.Get("addr")).
		Str("path", mpath).
		Msg("Start Video Server")

	if m.Stream == nil {
		m.Stream = mjpeg.NewStream()
	}
	http.Handle(mpath, m.Stream)

	mjpgQ := make(chan []byte)

	// go func the command listener
	go func() {
		log.Info().Msg("Starting MJPEG server")
		var cmd TLV
		for {
			select {
			case buf := <-mjpgQ:
				m.Stream.UpdateJPEG(buf)
				log.Warn().Str("cmd", cmd.Str()).Str("TODO", "implement this command")
			}
		}
	}()

	// Now go func the MJPEG HTTP server
	go http.ListenAndServe(config.Get("video-addr"), nil)
	return mjpgQ
}
