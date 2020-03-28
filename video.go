package main

import (
	"github.com/apex/log"
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

// VidePlayer will open and take control a single camera. At
// the moment any camera or device string that can be read
// by OpenCV are supported. A version
type VideoPlayer struct {

	// Check if we are recording
	recording     bool // XXX: mutex or channel this bool
	doAI          bool // TODO
	*mjpeg.Stream      // Stream will always be available
}

// NewVideoPlayer will create a new video player with default nil set.
func NewVideoPlayer(config *Configuration) (vid *VideoPlayer) {
	vid = &VideoPlayer{} // defaults are all good
	return vid
}

// Start Video opens the camera (sensor) and data (vidoe) starts streaming in.
// We will be streaming MJPEG for our initial use case.
func (vid *VideoPlayer) StartVideo() {
	var err error
	var buf []byte

	var cstr interface{}
	cstr = config.Camstr
	if cstr == "0" {
		cstr = 0
	}

	defer l.WithField("devid", cstr).Trace("entered start vid").Stop(&err)
	if vid.recording {
		l.Error("camera already recording")
		return
	}

	faced := NewFaceDetector()

	// Both API REST server and MQTT server have started up, we are
	// now waiting for requests to come in and instruct us wat to do.
	for img := range vid.StreamVideo(config.Camstr) {

		// Here we run through the AI, or whatever filter chain we are going
		// to use. For now it is hard coded with face detect, this will become
		// more flexible by allowing serial and concurrent filters.
		if vid.doAI {
			faced.FindFace(img)
		}

		// TODO: replace following when GoCV is not available.
		// Finalize the annotated image. XXX maybe we create a write channel?
		buf, err = gocv.IMEncode(".jpg", *img)
		if err != nil {
			l.Fatal("Failed encoding jpg")
		}
		vid.Stream.UpdateJPEG(buf)
	}
}

// StopVideo shuts the sensor down and turns
func (vid *VideoPlayer) StopVideo() {
	defer l.WithFields(log.Fields{
		"cameraid":  config.Camstr,
		"recording": vid.recording,
	}).Trace("Stop StreamVideo").Stop(nil)

	// Need to sync around this recording video (or can we use a channel)
	vid.recording = false
}

// StreamVideo takes a device string, starts the video stream and
// starts pumping single JPEG images from the camera stream.
func (vid *VideoPlayer) StreamVideo(devstr string) (frames <-chan *gocv.Mat) {
	var err error

	// Do not try to restart the video when it is already running.
	if vid.recording {
		l.Error("camera already recording")
		return nil
	}

	// Create the channel we are going to pump frames through
	frameQ := make(chan *gocv.Mat)

	defer l.WithFields(log.Fields{
		"camera":    devstr,
		"recording": vid.recording,
	}).Trace("StreamVideo").Stop(&err)

	// go function opens the webcam and starts reading from device, coping frames
	// to the frameQ processing channel
	go func() {
		var cam *gocv.VideoCapture

		camstr := GetCamstr(config.Camstr)

		log.Infof("Opening VideoCapture %s", camstr)

		// straight up 0
		//cam, err = gocv.OpenVideoCapture(camstr)
		cam, err = gocv.OpenVideoCapture(0)
		if err != nil {
			l.Fatal("failed to open video capture device")
			return
		}
		defer cam.Close()

		l.Info("Camera streaming  ...")

		// Only a single static image will be in the system at a given time.
		img := gocv.NewMat()

		// as long as vid.recording is true we will capture images and send
		// them into the image pipeline. We may recieve a REST or MQTT request
		// to stop recording, in that case the vid.recording will be set to
		// false and the recording will stop.
		vid.recording = true
		for vid.recording {

			// read a single raw image from the cam.
			if ok := cam.Read(&img); !ok {
				l.Info("device closed, turn recording off")
				vid.recording = false
			}

			// if the image is empty, there will be no sense continueing
			if img.Empty() {
				continue
			}

			// sent the frame to the frame pipeline
			frameQ <- &img
		}
	}()

	// return the frame channel, our caller will pass it to the reader
	return frameQ
}

// GetCamstr returns a string that OpenCV understands depending on the
// platform we are running on.
func GetCamstr(name string) (camstr string) {
	var ex bool
	if camstr, ex = camstrmap[name]; !ex {
		log.Errorf("camstr name %s NOT Found", name)
	}
	return camstr
}
