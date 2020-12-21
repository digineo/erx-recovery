package main

import (
	"log"
	"net"

	"github.com/vishvananda/netlink"
)

func getMACAddress(ifindex int, ipAddr net.IP) net.HardwareAddr {
	list, err := netlink.NeighList(ifindex, netlink.FAMILY_ALL)
	if err != nil {
		log.Fatalf("unable to get neighbor table of interface %v", localInterface.Name)
	}

	for i := range list {
		if list[i].IP.Equal(ipAddr) {
			return list[i].HardwareAddr
		}
	}

	return nil
}

func getInterfaceWithAddress(wanted net.IP) *net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for i := range ifaces {
		addrs, err := ifaces[i].Addrs()
		if err != nil {
			panic(err)
		}

		for j := range addrs {
			if addr := addrs[j].(*net.IPNet); addr != nil && addr.IP.Equal(wanted) {
				return &ifaces[i]
			}
		}
	}

	return nil
}
