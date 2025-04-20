package net

import (
	"fmt"
	"net"
)

func GetDefaultInterface() (*net.Interface, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80") // Google DNS
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP := localAddr.IP

	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iFace := range iFaces {
		addresses, err := iFace.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addresses {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.Equal(localIP) {
				return &iFace, nil
			}
		}
	}

	return nil, fmt.Errorf("no suitable interface was found")
}
