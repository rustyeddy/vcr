package main

import (
	"net"
	"os"

	"github.com/rs/zerolog/log"
)

func startupInfo() {
	if !config.Debug {
		return
	}
	log.Info().
		Str("app", "redeye").
		Str("pid", string(os.Getpid())).
		Str("hostname", GetHostname()).
		Msg("App is starting up ...")
}

// GetHostname for ourselves
func GetHostname() (hname string) {
	var err error
	if hname, err = os.Hostname(); err != nil {
		log.Error().Str("error", err.Error()).Msg("Good bye cruel world!")
	}
	return hname
}

// GetIPAddr of ourselves
func GetIPAddr() (addr string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error().Str("addr", addr).Str("error", err.Error()).Msg("Failed to get IP address")
		return
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
