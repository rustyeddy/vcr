package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/hybridgroup/mjpeg"
	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

var (
	video *VideoPlayer
)

// VidePlayer will open and take control a single camera. At
// the moment any camera or device string that can be read
// by OpenCV are supported. A version
type VideoPlayer struct {
	Name   string
	Addr   string
	Camstr string // String representing the camera

	// Video stream and a bool if we are recording
	*mjpeg.Stream `json:"-"` // Stream will always be available
	Recording     bool       `json:"recording"` // XXX: mutex or channel this bool

	// VideoPipeline filtering video. If nil, we have no filter or pipeline.
	VideoPipeline `json:"-"`
	SnapRequest   bool

	// Storage filename or directory
	Filename string
}

// GetVideoPlayer will create or return the video player
func GetVideoPlayer(config *Configuration) *VideoPlayer {
	if video == nil {
		video = &VideoPlayer{
			Name:     GetHostname(),
			Addr:     GetIPAddr(),
			Filename: "img/thumbnail.jpg",
			Camstr:   config.Camstr,
		} // defaults are all good

		if config.Pipeline != "" {
			video.SetPipeline(config.Pipeline)
		}

		// If the messanger is running, subscribe to our control topic
		if m := GetMessanger(); m != nil {
			m.Subscribe(video.GetControlChannel())
		}
	}
	return video
}

// NewVideoPlayer will create a new video player with default nil set.
func StartVideo(wg *sync.WaitGroup, config *Configuration) {
	defer wg.Done()

	// Set the route for video
	vpath := "/mjpeg"
	log.Info().
		Str("address", config.VideoAddr).
		Str("path", vpath).
		Msg("Start Video Server")

	video = GetVideoPlayer(config)

	video.Stream = mjpeg.NewStream()
	http.Handle(vpath, video.Stream)
	http.ListenAndServe(config.VideoAddr, nil)
}

// GetChannel returns the unique channel name for this camera
func (vid *VideoPlayer) GetAnnouncement() string {
	return vid.Addr + ":" + vid.Name
}

// GetChannel returns the unique channel name for this camera
func (vid *VideoPlayer) GetControlChannel() string {
	return "camera/" + vid.Name
}

// SetPipeline to be a named pipeline
func (vid *VideoPlayer) SetPipeline(name string) (err error) {
	vid.VideoPipeline, err = GetPipeline(name)
	return err
}

// Start Video opens the camera (sensor) and data (vidoe) starts streaming in.
// We will be streaming MJPEG for our initial use case.
func (vid *VideoPlayer) StartVideo() {
	var err error
	var buf []byte

	log.Info().Msg("StartVideo Entered ... ")
	defer log.Info().Msg(" XXX video has Finished")
	if vid.Recording {
		log.Warn().Msg("camera is already recording")
		return
	}

	// Both API REST server and MQTT server have started up, we are
	// now waiting for requests to come in and instruct us wat to do.
	for img := range vid.PumpVideo() {

		// Filter images if a VideoPipeline has been setup
		if vid.VideoPipeline != nil {
			vid.VideoPipeline.Send(img)
		}

		// TODO: replace following when GoCV is not available.
		// Finalize the annotated image. XXX maybe we create a write channel?
		buf, err = gocv.IMEncode(".jpg", *img)
		if err != nil {
			log.Fatal().Msg("Failed encoding jpg")
		}

		vid.Stream.UpdateJPEG(buf)

		// Check to see if a nsapshot has been requested, if so then
		// take a snapshot. TODO put this in the video pipeline
		if vid.SnapRequest {
			fname := "pub/img/snapshot.jpg"
			// Create the store

			var ok bool
			if ok = gocv.IMWrite(fname, *img); !ok {
				log.Error().Str("filename", fname).Msg("Snapshot failed to save ")
			}

			fmt.Printf("snap requested %s .. ok: %v ", fname, ok)
			log.Info().Str("filename", fname).Msg("Snapshot saved")

			vid.SnapRequest = false
		}
	}
	log.Info().Msg("Stopping Video")
}

// StreamVideo takes a device string, starts the video stream and
// starts pumping single JPEG images from the camera stream.
//
// TODO: Change this to Camera and create an interface that
// is sufficient for video files and imagnes.
func (vid *VideoPlayer) PumpVideo() (frames <-chan *gocv.Mat) {
	var err error

	// Do not try to restart the video when it is already running.
	if vid.Recording {
		log.Error().Msg("camera already recording")
		return nil
	}

	// Create the channel we are going to pump frames through
	frameQ := make(chan *gocv.Mat)
	defer log.Info().
		Str("cameraid", vid.Camstr).
		Bool("recording", vid.Recording).
		Msg("Stop StreamVideo")

	// go function opens the webcam and starts reading from device, coping frames
	// to the frameQ processing channel
	go func() {

		// Open the camera (capture device)
		var cam *gocv.VideoCapture
		camstr := GetCamstr(vid.Camstr)
		defer log.Info().
			Str("camstr", camstr).
			Msg("Opening VideoCapture")

		cam, err = gocv.OpenVideoCapture(camstr)
		if err != nil {
			log.Fatal().Msg("failed to open video capture device")
			return
		}
		defer cam.Close()


		log.Info().Msg("Camera streaming  ...")

		// Only a single static image will be in the system at a given time.
		img := gocv.NewMat()

		// as long as vid.recording is true we will capture images and send
		// them into the image pipeline. We may recieve a REST or MQTT request
		// to stop recording, in that case the vid.recording will be set to
		// false and the recording will stop.
		vid.Recording = true
		for vid.Recording {

			// read a single raw image from the cam.
			if ok := cam.Read(&img); !ok {
				log.Info().Msg("device closed, turn recording off")
				vid.Recording = false
			}

			// if the image is empty, there will be no sense continueing
			if img.Empty() {
				continue
			}

			// sent the frame to the frame pipeline (should we send by )
			frameQ <- &img
		}
	}()

	// return the frame channel, our caller will pass it to the reader
	return frameQ
}

// StopVideo shuts the sensor down and turns
func (vid *VideoPlayer) StopVideo() {
	defer log.Info().
		Str("cameraid", config.Camstr).
		Bool("recording", vid.Recording).
		Msg("Stop StreamVideo")

	// Need to sync around this recording video (or can we use a channel)
	vid.Recording = false
}

// GetCamstr returns a string that OpenCV understands depending on the
// platform we are running on.
func GetCamstr(name string) (camstr string) {
	var ex bool
	if camstr, ex = camstrmap[name]; !ex {
		log.Error().Str("name", name).Msg("camstr NOT Found")
	}
	return camstr
}
