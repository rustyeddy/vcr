package main

import (
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

// FaceDetector aides in identifing faces
type FaceDetector struct {
	Start     *time.Time
	End       *time.Duration
	FaceCount int
	gocv.CascadeClassifier
}

var (
	// color for the rect when faces detected
	blue         = color.RGBA{0, 0, 255, 0}
	faceDetector = &FaceDetector{}
)

// NewFaceDetector creates a new face detector
func NewFaceDetector() (fd *FaceDetector) {
	fd = &FaceDetector{}
	return fd
}

// Setup the video pipeline, basically read the classifier based on
// the particular haarscascade code.
func (fd *FaceDetector) Setup() {
	// load classifier to recognize faces
	fd.CascadeClassifier = gocv.NewCascadeClassifier()
	if !fd.CascadeClassifier.Load(config.XMLFile) {
		l.WithField("xmlfile", config.XMLFile).Error("Error reading cascade file")
	}
}

// FaceDetector takes in an image and finds a Face.
func (fd *FaceDetector) Send(img *gocv.Mat) *gocv.Mat {

	// detect faces
	rects := fd.CascadeClassifier.DetectMultiScale(*img)
	faceDetector.FaceCount = len(rects)

	// draw a rectangle around each face on the original image,
	// along with text identifying as "Human"
	for _, r := range rects {
		gocv.Rectangle(img, r, blue, 3)

		size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		gocv.PutText(img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
	}
	return img
}
