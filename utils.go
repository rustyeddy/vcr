package redeye

import (
	"net"
	"os"
	"log"
)

// GetHostname for ourselves
func GetHostname() (hname string) {
	var err error
	if hname, err = os.Hostname(); err != nil {
		log.Println("error: ", err.Error());
		return ""
	}
	return hname
}

// GetIPAddr of ourselves
func GetIPAddr() (addr string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("addr: ", addr, ", ", "error: ", err.Error(), " Failed to get IP address")
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
