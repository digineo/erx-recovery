package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func getMACAddress(iface *net.Interface, ipAddr net.IP) (net.HardwareAddr, error) {
	list, err := netlink.NeighList(iface.Index, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("unable to get neighbor table of interface %v: %w", iface.Name, err)
	}

	for i := range list {
		if list[i].IP.Equal(ipAddr) {
			return list[i].HardwareAddr, nil
		}
	}

	return nil, errors.New("MAC address not found")
}
