package models

import (
	"net"
	"os"
)

func GetHostIP() string {
	host, _ := os.Hostname()

	addrs, _ := net.LookupHost(host)
	for _, addr := range addrs {
		return addr
	}

	return "127.0.0.1"
}

func GetLocalIP() []string {
	ips := []string{}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}
