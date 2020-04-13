package main

import (
	"plugin"
	"sync"

	"gocv.io/x/gocv"
)

// VideoPipeline is a series of pipes that accepts and image
// processes the image, then returns the processed image. That
// image may then be process by another step in the pipe, which
// may include writting to a file or
type VideoPipeline interface {
	Name() string
	Setup() error
	Send(*gocv.Mat) *gocv.Mat
}

// FrameDrain listens to a channel delivering video camera images,
// typically to observe something and perhaps perform some
// transformation.
type VideoPipe struct {
	name string

	Q       chan *gocv.Mat // recieving data
	Process func(img *gocv.Mat)
	Next    VideoPipeline // try using this!!
}

var (
	pipelines map[string]VideoPipeline
)

func init() {
	pipelines = make(map[string]VideoPipeline)
}

// NewVideoPipe will create a new image Q to recieve upstream
// images, if the Process method is not nil, then the frame will
// go through the corresponding processing.
func NewVideoPipe(name string, pipe VideoPipeline) (fw *VideoPipe) {
	fw = &VideoPipe{
		name: name,
		Q:    make(chan *gocv.Mat), // no buffers for now
		Next: pipe,
	}
	return fw
}

// Name returns the name of this pipeline
func (p *VideoPipe) Name() string {
	return p.name
}

// Setup run optional setup.
func Setup() {
	// This function must exist to fullfill the VideoPipeline contract
}

// Listen for incoming images and send them to the frame pipeline
func (fq *VideoPipe) Listen(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	l.WithField("Name", fq.Name).Info("begin listen loop")

	// loop around waiting for incoming frames
	for {

		l.WithField("Name", fq.Name).Info("in loop")
		// Wait for a new frame and check that it is ok
		img, ok := <-fq.Q
		if !ok {
			l.Error("Listen channel appears to be closed")
			return
		}

		// If the Process callback has been set it will be called
		// for the image.
		if fq.Process != nil {
			fq.Process(img)
		}

		// TODO must be tested (how about sending on channel) ...
		if fq.Next != nil {
			fq.Next.Send(img)
		}
	}
}

// Send and Frame to the existing q, then if next exists, send it to the next.
// next will need to turn into a queue .. (hmm a channel of channels?)
func (fq *VideoPipe) Send(img *gocv.Mat) *gocv.Mat {
	l.Debug("sending image")
	fq.Q <- img
	return img
}

// GetPipeline will return a VideoPipeline `name` if one exists
// in the videoPipeline.
func GetPipeline(fname string) (p VideoPipeline, err error) {
	var ex bool
	if p, ex = pipelines[fname]; ex {
		return p, nil
	}

	pl, err := plugin.Open(fname)
	if err != nil {
		l.WithError(err).Error("failed to open plugin")
		return nil, err
	}

	sym, err := pl.Lookup("Pipeline")
	if err != nil {
		l.WithError(err).Error("Find the Pipe")
		return nil, err
	}
	p = sym.(VideoPipeline)

	// Run setup
	p.Setup()

	pipelines[fname] = p
	return p, nil
}
