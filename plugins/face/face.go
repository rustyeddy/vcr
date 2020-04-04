package main

import (
	"errors"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type FacePipeline struct {
	gocv.CascadeClassifier
}

var (
	// color for the rect when faces detected
	Blue     = color.RGBA{0, 0, 255, 0}
	Pipeline = FacePipeline{}
)

func init() {
	Pipeline.LoadClassifier("data/haarcascade_upperbody.xml")
}

// LoadClassier accept the filename of a HaarCascade .xml file
func (f *FacePipeline) LoadClassifier(fname string) (err error) {
	// load classifier to recognize faces
	f.CascadeClassifier = gocv.NewCascadeClassifier()
	if !f.CascadeClassifier.Load(fname) {
		return errors.New("Error reading cascade file: " + fname)
	}
	return nil
}

//
// ---------- Satisfy the VideoPipeline Interface --------

// FaceDetector takes in an image and finds a Face.
func (f *FacePipeline) Send(img *gocv.Mat) *gocv.Mat {
	// detect faces
	rects := f.CascadeClassifier.DetectMultiScale(*img)

	// draw a rectangle around each face on the original image,
	// along with text identifying as "Human"
	for _, r := range rects {
		gocv.Rectangle(img, r, Blue, 3)

		size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		gocv.PutText(img, "Human", pt, gocv.FontHersheyPlain, 1.2, Blue, 2)
	}
	return img
}
