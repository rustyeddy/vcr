package main

import (
	"net"
	"os"

	"github.com/apex/log"
)

func startupInfo() {
	if !config.Debug {
		return
	}
	log.Infof("config %v\n", config)

	l.WithFields(log.Fields{
		"app":      "redeye",
		"pid":      os.Getpid(),
		"hostname": GetHostname(),
	}).Info("App is starting up ...")
}

// GetHostname for ourselves
func GetHostname() (hname string) {
	var err error
	if hname, err = os.Hostname(); err != nil {
		log.WithError(err).Fatal("Good bye cruel world!")
	}
	return hname
}

// GetIPAddr of ourselves
func GetIPAddr() (addr string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.WithField("addr", addr).Fatal(err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
