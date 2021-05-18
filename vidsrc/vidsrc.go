package vidsrc

import "github.com/redeyelab/redeye/img"

type Vidsrc interface {
	Play()
	Pause()
	Snapshot()
	PumpVideo(frames <-chan *img.Frame)
}
