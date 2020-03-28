package main

import (
	"fmt"
	"sync"

	"gocv.io/x/gocv"
)

// FrameDrain listens to a channel delivering video camera images,
// typically to observe something and perhaps perform some
// transformation.
type FrameQ struct {
	Name    string
	Q       chan *gocv.Mat // recieving data
	Process func(img *gocv.Mat)
	Next    *FrameQ
}

// NewFrameQ will create a new image Q to recieve upstream
// images, if the Process method is not nil, then the frame will
// go through the corresponding processing.
func NewFrameQ(name string, proc func(img *gocv.Mat)) (fw *FrameQ) {
	fw = &FrameQ{
		Name:    name,
		Q:       make(chan *gocv.Mat), // no buffers for now
		Process: proc,
	}
	return fw
}

// Listen for incoming images and send them to the frame pipeline
func (fq *FrameQ) Listen(done <-chan bool, wg *sync.WaitGroup) {
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
func (fq *FrameQ) Send(img *gocv.Mat, next *FrameQ) {
	fmt.Println("sending")
	fq.Q <- img
}
