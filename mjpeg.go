package redeye

import (
	"log"
	"net/http"
	"github.com/hybridgroup/mjpeg"
)

// MJPEGServer starts up an HTTP server for posting rapidly
// updating JPEG images from the camera.
type MJPEGServer struct {
	Name          string
	Addr          string
	*mjpeg.Stream `json:"-"` // Stream will always be available

	Q chan []byte
}

var (
	mjp *MJPEGServer
)

func GetMJPEGServer() *MJPEGServer {
	if mjp == nil {
		mjp = NewMJPEGServer()
	}
	return mjp
}

// NewMJPEGServer will create a new video player with default nil set.
func NewMJPEGServer() (m *MJPEGServer) {
	m = &MJPEGServer{"mjpg", ":8887", nil, make(chan []byte)}
	return m
}

// Start the mjpeg server clients will likely be either the web client
// or the CV pipeline.
func (m *MJPEGServer) Start(cmdQ chan TLV) (mpgQ chan []byte) {

	// Set the route for video
	mpath := "/mjpeg"
	if Config.Debug {
		log.Println(m.Addr, " path ", mpath, " Starting video server")		
	}

	if m.Stream == nil {
		m.Stream = mjpeg.NewStream()
	}
	http.Handle(mpath, m.Stream)

	mjpgQ := make(chan []byte)

	// go func the command listener
	go func() {
		if Config.Debug {
			log.Println("Starting MJPEG server")			
		}
	
		var cmd TLV
		for {
			select {
			case buf := <-mjpgQ:
				m.Stream.UpdateJPEG(buf)
				log.Println("TODO: implement this command: ", cmd.Str())
			}
		}
	}()

	// Now go func the MJPEG HTTP server: make this a config
	addr := "http://localhost:8833/mjpeg"
	if Config.Debug {
		log.Println("MJPEG Streaming video addrress: ", addr)		
	}

	go http.ListenAndServe(addr, nil)
	return mjpgQ
}
