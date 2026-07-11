package main

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strconv"
	"strings"
)

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

func extractMAC(out []byte, ipAddr net.IP) (net.HardwareAddr, error) {
	ipStr := "(" + ipAddr.String() + ")"
	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) >= 4 && fields[1] == ipStr {
			return parseMAC(fields[3])
		}
	}
	return nil, errors.New("IP address not found")
}

func parseMAC(macStr string) (net.HardwareAddr, error) {
	fields := strings.Split(macStr, ":")
	if len(fields) != 6 {
		return nil, errors.New("invalid MAC address format")
	}

	mac := make(net.HardwareAddr, 6)

	for i, field := range fields {
		if len(field) < 1 || len(field) > 2 {
			return nil, errors.New("invalid MAC address format")
		}

		value, err := strconv.ParseUint(field, 16, 8)
		if err != nil {
			return nil, err
		}

		mac[i] = byte(value)
	}

	return mac, nil
}
