//go:build !linux

package main

import (
	"log"
	"net"
	"os/exec"
)

func getMACAddress(iface *net.Interface, ipAddr net.IP) (net.HardwareAddr, error) {
	out, err := exec.Command("arp", "-an", "-i", iface.Name).Output()
	if err != nil {
		log.Fatal(err)
	}

	return extractMAC(out, ipAddr)
}
