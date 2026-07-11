//go:build !linux

package main

import (
	"bytes"
	"net"
	"os"
	"testing"
)

func TestGetMACAddress(t *testing.T) {
	out, err := os.ReadFile("testdata/arp")
	if err != nil {
		t.Fatal(err)
	}

	extractedMAC, err := extractMAC(out, net.ParseIP("192.168.1.1"))
	if err.Error() != "IP address not found" {
		t.Fatal(err)
	}

	extractedMAC, err = extractMAC(out, net.ParseIP("192.168.180.10"))
	if err != nil {
		t.Fatal(err)
	}

	expectedMAC, err := parseMAC("52:54:0:ff:15:70")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(extractedMAC, expectedMAC) {
		t.Errorf("expected MAC %s, got %s", expectedMAC, extractedMAC)
	}
}
func TestParseMACAddress(t *testing.T) {
	mac, err := parseMAC("52:54:0:ff:b:70")
	if err != nil {
		t.Fatal(err)
	}
	if mac.String() != "52:54:00:ff:0b:70" {
		t.Errorf("expected MAC 52:54:00:ff:0b:70, got %s", mac.String())
	}
}
