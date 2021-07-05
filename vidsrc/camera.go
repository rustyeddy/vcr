package vidsrc

import (
	"fmt"
	"log"

	"gocv.io/x/gocv"
)

type Camera struct {
	Camstr      string
	Recording   bool
	SnapRequest bool
}

func NewCamera(camstr string) (cam *Camera) {
	cam = &Camera{
		Recording: false,
		Camstr:    camstr,
	}
	return cam
}

func (cam *Camera) Play() {
	cam.Recording = true
}

func (cam *Camera) Pause() {
	cam.Recording = false
}

func (cam *Camera) Snap() {
	cam.Recording = true
}

// StreamVideo takes a device string, starts the video stream and
// starts pumping single JPEG images from the camera stream.
//
// TODO: Change this to Camera and create an interface that
// is sufficient for video files and imagnes.
//func (cam *Camera) PumpVideo() (frames <-chan *gocv.Mat) {
func (cam *Camera) PumpVideo() (frames <-chan *gocv.Mat) {
// func (cam *Camera) PumpVideo() (frames <-chan *img.Frame) {

	// Do not try to restart the video when it is already running.
	if cam.Recording {
		log.Println("camera already recording")
		return nil
	}

	// Create the channel we are going to pump frames through
	frameQ := make(chan *gocv.Mat)
	// frameQ := make(chan *img.Frame)

	// go function opens the webcam and starts reading from device, coping frames
	// to the frameQ processing channel
	go func() {

		defer log.Println("cameraid: ", cam.Camstr, " recording: ", cam.Recording, " Stop StreamVideo")

		// // Open the camera (capture device)
		var cap *gocv.VideoCapture
		camstr := GetCamstr(cam.Camstr)
		log.Println("Opening Video with camstr: ", camstr, "Opening VideoCapture")

		var err error
		cap, err = gocv.OpenVideoCapture(camstr)
		if err != nil {
			log.Fatal("failed to open video capture device")
			return
		}
		defer cap.Close()

		log.Println("Camera streaming  ...")

		// Only a single static image will be in the system at a given time.
		img := gocv.NewMat()

		// as long as cam.recording is true we will capture images and send
		// them into the image pipeline. We may recieve a REST or MQTT request
		// to stop recording, in that case the cam.recording will be set to
		// false and the recording will stop.
		cam.Recording = true
		for cam.Recording {

			// read a single raw image from the cam.
			if ok := cap.Read(&img); !ok {
				log.Println("device closed, turn recording off")
				cam.Recording = false
			}

			// if the image is empty, there will be no sense continueing
			if img.Empty() {
				continue
			}

			// sent the frame to the frame pipeline (should we send by )
			fmt.Printf("frame %+v\n", img)
			frameQ <- &img
		}
		log.Println("Recording: ", cam.Recording, " Video loop exiting ...")
	}()

	// return the frame channel, our caller will pass it to the reader
	return frameQ
}
