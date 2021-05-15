package main

import (

	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

var (
	video *VideoPlayer
)

// VidePlayer will open and take control a single camera. At
// the moment any camera or device string that can be read
// by OpenCV are supported. A version
type VideoPlayer struct {
	camera.Camera

	// VideoPipeline filtering video. If nil, we have no filter or pipeline.
	VideoPipeline `json:"-"`

	// Channel to send video on
	Q chan TLV
}

// GetVideoPlayer
func GetVideoPlayer() (v *VideoPlayer) {
	if video == nil {
		v = NewVideoPlayer()
	}
	return v
}

// GetVideoPlayer will create or return the video player.
// TODO Change this to accept a configmap
func NewVideoPlayer() (video *VideoPlayer) {
	video = &VideoPlayer{}
	video.Camstr = config.Vidsrc

	if config.Pipeline != "" {
		video.SetPipeline(config.Pipeline)
	}
	return video
}

// NewVideoPlayer will create a new video player with default nil set.
func (vid *VideoPlayer) Start(cmdQ chan TLV) chan TLV {

	// go func the command listener
	go func() {
		log.Info().Msg("Starting Video service listener .. ")
		for {
			select {
			case cmd := <-cmdQ:
				log.Info().Str("cmd", cmd.Str()).Msg("incoming video command")
				switch cmd.Type() {
				case CMDPlay:
					log.Info().Str("cmd", "play").Msg("Playing Video...")
					go vid.Play()

				case CMDPause:
					log.Info().Str("cmd", "pause").Msg("Pausing Video...")
					vid.Pause()

				default:
					log.Warn().Str("cmd", cmd.Str()).Msg("unknown command")
				}
			}
		}
	}()
	return vid.Q
}

// SetPipeline to be a named pipeline
func (vid *VideoPlayer) SetPipeline(name string) (err error) {
	vid.VideoPipeline, err = GetPipeline(name)
	return err
}

// Start Video opens the camera (sensor) and data (vidoe) starts streaming in.
// We will be streaming MJPEG for our initial use case.
func (vid *VideoPlayer) Play() {
	var err error
	var buf []byte

	log.Info().Msg("StartVideo Entered ... ")
	defer log.Info().Msg("StartVideo video has exited")
	if vid.Recording {
		log.Warn().Msg("camera is already recording")
		return
	}

	// Both API REST server and MQTT server have started up, we are
	// now waiting for requests to come in and instruct us wat to do.
	for img := range vid.PumpVideo() {

		// Filter images if a VideoPipeline has been setup
		if vid.VideoPipeline != nil {
			vid.VideoPipeline.Send(img)
		}

		// TODO: replace following when GoCV is not available.
		// Finalize the annotated image. XXX maybe we create a write channel?
		buf, err = gocv.IMEncode(".jpg", *img)
		if err != nil {
			log.Fatal().Str("comp", "video").Msg("Failed encoding jpg")
		}

		mjp := GetMJPEGServer()

		// Send the annotated buffer to the MJPEG server
		mjp.Q <- buf

		// Check to see if a nsapshot has been requested, if so then
		// take a snapshot. TODO put this in the video pipeline
		if vid.SnapRequest {
			fname := "pub/img/snapshot.jpg"
			// Create the store

			var ok bool
			if ok = gocv.IMWrite(fname, *img); !ok {
				log.Error().Str("filename", fname).Msg("Snapshot failed to save ")
			}
			log.Info().Str("filename", fname).Msg("Snapshot saved")
			vid.SnapRequest = false
		}
	}
	log.Info().Msg("Stopping Video")
}

// StopVideo shuts the sensor down and turns
func (vid *VideoPlayer) Pause() {
	defer log.Info().
		Str("cameraid", vid.Camstr).
		Bool("recording", vid.Recording).
		Msg("Stop StreamVideo")

	// Need to sync around this recording video (or can we use a channel)
	vid.Recording = false
}

// VideoPlayerStatus is returned by the REST api reporting
// the status of
type VideoPlayerStatus struct {
	Addr      string
	Camstr    string
	Recording bool
	Pipeline  string
}

func (vid *VideoPlayer) Status() (status *VideoPlayerStatus) {
	status = &VideoPlayerStatus{
		Camstr:    vid.Camstr,
		Recording: vid.Recording,
	}
	if vid.VideoPipeline != nil {
		status.Pipeline = vid.VideoPipeline.Name()
	}
	return status
}
