package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"

	"github.com/google/goterm/term"
	"github.com/pin/tftp"
)

const kernelFilename = "vme50"

func startTFTP() {
	bindAddr := net.JoinHostPort(localAddr.String(), "69")
	s := tftp.NewServer(tftpReadHandler, nil)

	go func() {
		log.Println("listening on", bindAddr)
		err := s.ListenAndServe(bindAddr) // blocks until s.Shutdown() is called
		if err != nil {
			fmt.Fprintf(os.Stdout, "server: %v\n", err)
			os.Exit(1)
		}
	}()
}

// tftpReadHandler is called when client starts file download from server
func tftpReadHandler(filename string, rf io.ReaderFrom) error {
	if filename != kernelFilename {
		fmt.Fprintf(os.Stderr, "filename mismatch, requested=%s expected=%s", filename, kernelFilename)
		return os.ErrNotExist
	}

	remoteAddr := rf.(tftp.OutgoingTransfer).RemoteAddr().IP

	// Figure out remote MAC address
	mac := getMACAddress(localInterface.Index, remoteAddr)
	if mac == nil {
		log.Fatalf("unable to get MAC address for %v", remoteAddr)
	}

	log.Println(term.Greenf("Download started from %v (%v)", remoteAddr, mac))

	// Set the logfile path
	writer.setFile(path.Join(*logFolder, fmt.Sprintf("%x.log", []byte(mac))))

	// open image file
	file, err := os.Open(*kernelPath)
	if err != nil {
		log.Println(term.Red(err.Error()))
		return err
	}

	// read image file
	n, err := rf.ReadFrom(file)
	if err != nil {
		log.Println(term.Red(err.Error()))
		return err
	}

	log.Println(term.Greenf("Download finished (%d bytes sent)", n))
	return nil
}
