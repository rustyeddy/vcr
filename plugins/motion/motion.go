package main

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

const MinimumArea = 3000

type MotionPipeline struct {
	name string
}

var (
	Pipeline = MotionPipeline{"motion"}

	imgDelta  = gocv.NewMat()
	imgThresh = gocv.NewMat()
	mog2      = gocv.NewBackgroundSubtractorMOG2()

	status = "Ready"

	// color for the rect when faces detected
	Red   = color.RGBA{255, 0, 0, 0}
	Green = color.RGBA{0, 255, 0, 0}
	Blue  = color.RGBA{0, 0, 255, 0}
)

func (m *MotionPipeline) Name() string {
	return m.name
}

func (m *MotionPipeline) Setup() error {
	return nil
}

func (m *MotionPipeline) Send(img *gocv.Mat) *gocv.Mat {
	status = "Ready"
	statusColor := Blue

	// Objtain foreground image
	mog2.Apply(*img, &imgDelta)

	// remaining cleanup of the image to use for finding contours
	// first use the threshold
	gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

	// now dilate
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	defer kernel.Close()

	gocv.Dilate(imgThresh, &imgThresh, kernel)

	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	for i, c := range contours {
		area := gocv.ContourArea(c)
		if area < MinimumArea {
			continue
		}

		status = "Motion Detected"
		statusColor = Red
		gocv.DrawContours(img, contours, i, statusColor, 2)

		rect := gocv.BoundingRect(c)
		gocv.Rectangle(img, rect, Green, 2)
	}
	gocv.PutText(img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

	return img
}
