package main

import (
	"errors"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type FacePipeline struct {
	name    string
	xmlFile string

	gocv.CascadeClassifier
}

var (
	// color for the rect when faces detected
	Blue     = color.RGBA{0, 0, 255, 0}
	Pipeline = FacePipeline{
		name:    "face",
		xmlFile: "data/haarcascade_upperbody.xml",
	}
)

func init() {
}

// Name is the name of the pipe line
func (f *FacePipeline) Name() string {
	return f.name
}

// Setup allows us to setup the plugins
func (f *FacePipeline) Setup() error {
	err := f.LoadClassifier()
	if err != nil {
		return err
	}
	return nil
}

// LoadClassier accept the filename of a HaarCascade .xml file
func (f *FacePipeline) LoadClassifier() (err error) {
	// load classifier to recognize faces
	f.CascadeClassifier = gocv.NewCascadeClassifier()
	if !f.CascadeClassifier.Load(f.xmlFile) {
		return errors.New("Error reading cascade file: " + f.xmlFile)
	}
	return nil
}

//
// ---------- Satisfy the VideoPipeline Interface --------

// FaceDetector takes in an image and finds a Face.
func (f *FacePipeline) Send(img *gocv.Mat) *gocv.Mat {
	// detect facesa

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
