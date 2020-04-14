package main

import (
	"flag"
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
	Output       string `json:"output"`

	Loglevel string

	MQTT string `json:"mqtt"`

	Pipeline   string `json:"pipeline"` // name of the plugin
	XMLFile    string `json:"xmlfile"`
}

var (
	config Configuration
)

func init() {
	flag.StringVar(&config.Addr, "address", "0.0.0.0:8888", "web address default 0.0.0.0:8888")
	flag.StringVar(&config.MQTT, "mqtt", "tcp://10.24.10.10:1883", "mqtt broker address def tcp://10.24.10.10:1883")
	flag.StringVar(&config.StaticPath, "pub", "./pub", "Application root dir")
	flag.StringVar(&config.IndexPath, "index", "index.html", "index file")
	flag.StringVar(&config.Name, "name", "redeye", "Application Name")
	flag.StringVar(&config.Camstr, "camstr", "0", "Camera ID")
	flag.StringVar(&config.ConfigFile, "config", "redeye.json", "Config file: redeye.json")
	flag.StringVar(&config.Loglevel, "loglevel", "info", "default log level is debug")
	flag.StringVar(&config.VideoAddr, "video-addr", "0.0.0.0:8887", "web address default 0.0.0.0:8887")
	flag.StringVar(&config.Pipeline, "pipeline", "", "Face detect")

	flag.BoolVar(&config.ServeVideo, "serve-video", true, "display video on local screen if available")
	flag.BoolVar(&config.DisplayVideo, "display-video", true, "display video on local screen if available")

	flag.StringVar(&config.XMLFile, "xmlfile", "data/haarcascade_frontalface_alt2.xml", "XMLFile")
}
