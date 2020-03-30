package main

import (
	"fmt"
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

// FaceDetector takes in an image and finds a Face.
func (fd *FaceDetector) Send(img *gocv.Mat) *gocv.Mat {

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(config.XMLFile) {
		fmt.Printf("Error reading cascade file: %v\n", config.XMLFile)
		return nil
	}

	// detect faces
	rects := classifier.DetectMultiScale(*img)
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
