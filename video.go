package redeye

import (
	"fmt"
	"log"

	"github.com/redeyelab/redeye/vidsrc"
)

var (
	video *VideoPlayer
)

// VidePlayer will open and take control a single camera. At
// the moment any camera or device string that can be read
// by OpenCV are supported. A version
type VideoPlayer struct {
	vidsrc.Camera          // where the videos come from change to generic
	Q             chan TLV // where to send the video

	VideoPipeline
}

// GetVideoPlayer
func GetVideoPlayer() (v *VideoPlayer) {
	if video == nil {
		v = NewVideoPlayer()
	}
	return v
}

// GetVideoPlayer will create or return the video player.
// TODO Change this to accept a configmap
func NewVideoPlayer() (video *VideoPlayer) {
	video = &VideoPlayer{}
	video.Camstr = "/dev/video1"
	return video
}

// NewVideoPlayer will create a new video player with default nil set.
func (vid *VideoPlayer) Start(cmdQ chan TLV) chan TLV {

	// go func the command listener
	go func() {
		if Config.Debug {
			log.Println("Starting Video service listener .. ")
		}
		for {
			select {
			case cmd := <-cmdQ:
				if Config.Debug {
					log.Println("incoming video command")
				}
				switch cmd.Type() {
				case CMDPlay:
					go vid.Play()

				case CMDPause:
					vid.Pause()

				default:
					log.Printf("VidPlayer Start: unknown command %+v\n", cmd)
				}
			}
		}
	}()
	return vid.Q
}

// SetPipeline to be a named pipeline
func (vid *VideoPlayer) SetPipeline(name string) (err error) {
	vid.VideoPipeline, err = GetPipeline(name)
	return err
}

// Start Video opens the camera (sensor) and data (vidoe) starts streaming in.
// We will be streaming MJPEG for our initial use case.
func (vid *VideoPlayer) Play() error {
	var err error
	var buf []byte

	log.Println("StartVideo Entered ... ")
	defer log.Println("StartVideo video has exited")
	if vid.Recording {
		return fmt.Errorf("camera is already recording")
	}

	// Both API REST server and MQTT server have started up, we are
	// now waiting for requests to come in and instruct us wat to do.
	for frame := range vid.PumpVideo() {

		// Filter images if a VideoPipeline has been setup
		if vid.VideoPipeline != nil {
			vid.VideoPipeline.Send(frame)
		}

		// TODO: replace following when GoCV is not available.
		// Finalize the annotated image. XXX maybe we create a write channel?
		// buf, err = gocv.IMEncode(".jpg", *frame)
		// if err != nil {
		//	log.Fatal().Str("comp", "video").Msg("Failed encoding jpg")
		//}

		mjp := GetMJPEGServer()

		// Send the annotated buffer to the MJPEG server
		mjp.Q <- buf

		// Check to see if a nsapshot has been requested, if so then
		// take a snapshot. TODO put this in the video pipeline
		if vid.SnapRequest {
			fname := "pub/img/snapshot.jpg"
			// Create the store

			if err = frame.Save(fname); err != nil {
				return fmt.Errorf("filename: snapshot save failed %s", fname)
			}
			vid.SnapRequest = false
		}
	}
	return nil
}

// StopVideo shuts the sensor down and turns
func (vid *VideoPlayer) Pause() {
	if Config.Debug {
		defer log.Println("camera-id: ", vid.Camstr, " recording: ", " Stop StreamVideo")
	}
	vid.Recording = false
}

// VideoPlayerStatus is returned by the REST api reporting
// the status of
type VideoPlayerStatus struct {
	Addr      string
	Camstr    string
	Recording bool
	Pipeline  string
}

func (vid *VideoPlayer) Status() (status *VideoPlayerStatus) {
	status = &VideoPlayerStatus{
		Camstr:    vid.Camstr,
		Recording: vid.Recording,
	}
	if vid.VideoPipeline != nil {
		status.Pipeline = vid.VideoPipeline.Name()
	}
	return status
}
