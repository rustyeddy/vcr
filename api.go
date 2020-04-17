/*

RedEye the smart camera software.

- MQTT Play and Pause Video

- GET		/health
- GET		/config
- POST|PUT	/config/?key=val&key=val
- GET		/messanger

*/
package redeye

// Frame is a wrapper around the source of the frame. Specifically
// this interface should be able to define a pretty light weight
// video support, unfortunantly OpenCV, as **wonderful** has it is
// is not exactly light weight.  Support for v4l, or other
// frameworks should be possible without a forklift change
type Frame interface {
	Raw() []byte
	Len() int

	// XXX I think i need these?
	Reader(p []byte) (n int, err error)
	Writer(p []byte) (n int, err error)
}

// VideoSource includes cameras, files, the network and
// still images, like jpgs, etc.
type VideoSource interface {
	Config()
	Start() (frameQ chan Frame)
	Stop()
}

// VideoSink recieves a video stream, then prepares and sends
// the frames to their ultimate destination, which typically
// includes the video server a display or storage of some form.
type VideoSink interface {
	ReadStream(frameQ chan Frame)
	Write(b []byte) (n int, err error)
}

// Pipe accepts a Frame as input, examines, modifies or something
// then returns a frame, possibly modified by the previous operation
type Pipe interface {
	Name() string
	Process(q chan Frame) *Frame
	Next() *Pipe
}

// Pipeline connects a
type Pipeline interface {
	SetVideoSource(vsrc VideoSource)
	SetVideoSink(vdst VideoSink)
	Append(p Pipe)
}

// Configmap
type Configmap interface {
	Name() string
	Keys() []string
	Exists(key string) bool
	Get(key string) (val string)
	Set(key string, val string) error
	Del(key string) error

	// Some helpers
	Int(key string) (val int)
	Str(key string) (val string)
}

// CameraStatus is passed along in the REST call
type CameraStatus struct {
	Name     string
	Addr     string
	Status   string
	Pipeline string
}

type Pipelines struct {
	Pipes []string
}
