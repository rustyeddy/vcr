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
	Name string
	Addr string

	// Video stream and a bool if we are recording
	*mjpeg.Stream `json:"-"` // Stream will always be available
	Recording     bool       `json:"recording"` // XXX: mutex or channel this bool

	// VideoPipeline filtering video. If nil, we have no filter or pipeline.
	VideoPipeline `json:"-"`
}

// NewVideoPlayer will create a new video player with default nil set.
func NewVideoPlayer(config *Configuration) (vid *VideoPlayer) {
	vid = &VideoPlayer{
		Name: GetHostname(),
		Addr: GetIPAddr(),
	} // defaults are all good

	if config.Pipeline != "" {
		vid.SetPipeline(config.Pipeline)
	}

	return vid
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

	l.Info("StartVideo Entered ... ")
	defer l.Info("StartVideo Finished")

	// This is pretty simple, almost every system that openCV supports
	// will use an integer or a string, for example I have testing this
	// on the following devices with the respective strings
	//
	// ubuntu-amd64 /dev/video0 v4l
	// macos 0 builtin
	// raspberry-pi 0 CSI
	// nano pipeline gstreamer-pipeline
	//
	// use -camstr to make sure it comes out correctly
	var cstr interface{}
	cstr = config.Camstr
	if cstr == "0" {
		cstr = 0
	}

	defer l.WithField("devid", cstr).Trace("entered start vid").Stop(&err)
	if vid.Recording {
		l.Error("camera already recording")
		return
	}

	// Video pipeline are named. Setting them is as simple as passing
	// in the name.
	//vid.SetPipeline("face")

	// Both API REST server and MQTT server have started up, we are
	// now waiting for requests to come in and instruct us wat to do.
	for img := range vid.StreamVideo(config.Camstr) {

		// Here we run through the AI, or whatever filter chain we are going
		// to use. For now it is hard coded with face detect, this will become
		// more flexible by allowing serial and concurrent filters.

		//if vid.doAI {
		//	faced.FindFace(img)
		//}
		if vid.VideoPipeline != nil {
			vid.VideoPipeline.Send(img)
		}

		// TODO: replace following when GoCV is not available.
		// Finalize the annotated image. XXX maybe we create a write channel?
		buf, err = gocv.IMEncode(".jpg", *img)
		if err != nil {
			l.Fatal("Failed encoding jpg")
		}
		vid.Stream.UpdateJPEG(buf)
	}
	l.Info("")
}

// StopVideo shuts the sensor down and turns
func (vid *VideoPlayer) StopVideo() {
	defer l.WithFields(log.Fields{
		"cameraid":  config.Camstr,
		"recording": vid.Recording,
	}).Trace("Stop StreamVideo").Stop(nil)

	// Need to sync around this recording video (or can we use a channel)
	vid.Recording = false
}

// StreamVideo takes a device string, starts the video stream and
// starts pumping single JPEG images from the camera stream.
func (vid *VideoPlayer) StreamVideo(devstr string) (frames <-chan *gocv.Mat) {
	var err error

	// Do not try to restart the video when it is already running.
	if vid.Recording {
		l.Error("camera already recording")
		return nil
	}

	// Create the channel we are going to pump frames through
	frameQ := make(chan *gocv.Mat)

	defer l.WithFields(log.Fields{
		"camera":    devstr,
		"recording": vid.Recording,
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
		vid.Recording = true
		for vid.Recording {

			// read a single raw image from the cam.
			if ok := cam.Read(&img); !ok {
				l.Info("device closed, turn recording off")
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

// GetCamstr returns a string that OpenCV understands depending on the
// platform we are running on.
func GetCamstr(name string) (camstr string) {
	var ex bool
	if camstr, ex = camstrmap[name]; !ex {
		log.Errorf("camstr name %s NOT Found", name)
	}
	return camstr
}
