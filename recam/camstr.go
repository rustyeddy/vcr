package main

// The cam string is a map of strings that can be used by OpenCV to
// open up and start recording from a camera.

var camstrmap map[string]string

// Make some hard coded decisions for now. We will test
func init() {
	camstrmap = map[string]string{
		"jetson":  jetsonCamstr(),
		"nano":    jetsonCamstr(),
		"rpi":     "0",
		"mac":     "0",
		"linux":   "/dev/video0",
		"default": "0",
		"0":       "0",
	}
}

func jetsonCamstr() string {
	gstpipe := "nvarguscamerasrc ! " +
		"video/x-raw(memory:NVMM), width=(int)1280, height=(int)720, format=(string)NV12, framerate=(fraction)60/1 ! " +
		"nvvidconv flip-method=0 ! " +
		"video/x-raw, width=(int)1280, height=(int)720, format=(string)BGRx ! " +
		"videoconvert ! " +
		"video/x-raw, format=(string)BGR !" +
		"appsink	"
	return gstpipe
}
