package main

import (
	"flag"

	"github.com/apex/log"
)

// Configuration structure for our camera software
type Configuration struct {
	Name string `json:"name"`
	Addr string `json:"addr"`

	StaticPath string `json:"static_path"` // root of the website
	IndexPath  string `json:"index-path"`  // name of default index file

	ConfigFile string `json:"-"` // do not save the name of the file
	Debug      bool   `json:"-"`

	VideoAddr string `json:"video"`  // Address this video is available on
	Camstr    string `json:"camstr"` // string to fireup the web cam

	DisplayVideo bool   `json:"display"`
	ServeVideo   bool   `json:"video"`
	FaceDetect   bool   `json:"face-detect"`
	XMLFile      string `json:"xmlfile"`
	Output       string `json:"output"`

	Loglevel string

	MQTT string `json:"mqtt"`
}

// GetConfig returns the one true Configuration with defaults set and command line
// flags ready to fly.
func GetConfig() *Configuration {
	var c Configuration
	flag.StringVar(&c.Addr, "address", "0.0.0.0:8888", "web address default 0.0.0.0:8888")
	flag.StringVar(&c.MQTT, "mqtt", "tcp://10.24.10.10:1883", "mqtt broker address def tcp://10.24.10.10:1883")
	flag.StringVar(&c.StaticPath, "pub", "./pub", "Application root dir")
	flag.StringVar(&c.IndexPath, "index", "index.html", "index file")
	flag.StringVar(&c.Name, "name", "redeye", "Application Name")
	flag.StringVar(&c.Camstr, "camstr", "0", "Camera ID")
	flag.StringVar(&c.ConfigFile, "config", "redeye.json", "Config file: redeye.json")
	flag.StringVar(&c.Loglevel, "loglevel", "info", "default log level is debug")

	flag.StringVar(&c.VideoAddr, "video-addr", "0.0.0.0:8887", "web address default 0.0.0.0:8887")

	flag.BoolVar(&c.ServeVideo, "serve-video", true, "display video on local screen if available")
	flag.BoolVar(&c.DisplayVideo, "display-video", true, "display video on local screen if available")

	flag.BoolVar(&c.FaceDetect, "face-detect", true, "run face detection algorithm")
	flag.StringVar(&c.XMLFile, "xmlfile", "data/haarcascade_frontalface_alt2.xml", "XMLFile")

	return &c
}

// ReadConfig reads the Configuration from a file.
func ReadConfig(name string) (c *Configuration) {
	log.Fatal("TODO: Write the configuration file")

	return c
}

// SaveConfig file will save the Configuration in a file.
func SaveConfig(name string, c *Configuration) {
	log.Fatal("TODO: Read the configuration file")
}
