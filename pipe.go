package main

import (
	"sync"

	"gocv.io/x/gocv"
)

// VideoPipeline is a series of pipes that accepts and image
// processes it and spits the processed image out the other end
type VideoPipeline interface {
	Send(*gocv.Mat) *gocv.Mat
}

// FrameDrain listens to a channel delivering video camera images,
// typically to observe something and perhaps perform some
// transformation.
type VideoPipe struct {
	Q       chan *gocv.Mat // recieving data

	Name    string
	Process func(img *gocv.Mat)
	Next    VideoPipeline
}

// NewVideoPipe will create a new image Q to recieve upstream
// images, if the Process method is not nil, then the frame will
// go through the corresponding processing.
func NewVideoPipe(name string, pipe VideoPipeline) (fw *VideoPipe) {
	fw = &VideoPipe{
		Name:    name,
		Q:       make(chan *gocv.Mat), // no buffers for now
		Next: pipe,
	}
	return fw
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
	}
}

// Send and Frame to the existing q, then if next exists, send it to the next.
// next will need to turn into a queue .. (hmm a channel of channels?)
func (fq *VideoPipe) Send(img *gocv.Mat) *gocv.Mat {
	l.Debug("sending image")
	fq.Q <- img
	return img
}
