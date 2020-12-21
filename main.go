package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	expect "github.com/google/goexpect"
	"github.com/google/goterm/term"
)

var (
	deviceIP   = flag.String("device-ip", "172.16.3.212", "device IP address for the TFTP client, needs to be in the same /24 as the TFTP server")
	localIP    = flag.String("local-ip", "172.16.3.210", "local IP address for the TFTP server")
	kernelPath = flag.String("image", "/path/to/kernel", "path to linux kernel image")
	logFolder  = flag.String("log", "./logs", "directory for log files")
	ttyPath    = flag.String("tty", "/dev/ttyUSB0", "path to serial console")
	verbose    = flag.Bool("verbose", true, "Increase verbosity")

	filesize       int
	e              *expect.GExpect
	writer         logWriter
	localAddr      net.IP
	localInterface *net.Interface
)

func main() {
	flag.Parse()

	// Find local interface
	localAddr = net.ParseIP(*localIP)
	localInterface = getInterfaceWithAddress(localAddr)
	if localInterface == nil {
		log.Fatalf("unable to find local interface with IP address %v", *localIP)
	}
	log.Printf("using local inteface %v with address %v", localInterface.Name, *localIP)

	// Set log folder
	if _, err := os.Stat(*logFolder); os.IsNotExist(err) {
		os.Mkdir(*logFolder, 0o755)
	}

	if stat, err := os.Stat(*kernelPath); err != nil {
		panic(err)
	} else {
		filesize = int(stat.Size())
	}

	startTFTP()

	if err := run(); err == nil {
		log.Println(term.Green("Boot successful"))
	} else {
		log.Println(term.Red(err.Error()))
		os.Exit(1)
	}
}

func run() (err error) {
	e, _, err = expect.SpawnWithArgs([]string{"picocom", "-b", "57600", *ttyPath}, time.Second,
		expect.Tee(&writer),
		expect.Verbose(*verbose),
	)
	if err != nil {
		return err
	}
	defer e.Close()
	writer.Close()

	waitForUBoot()
	flashAndReboot()
	waitForUBoot()

	return verifyBoot()
}

var (
	statusFlashCorrupt = expect.NewStatus(0, "flash corrupt")
	reBadMagicNumber   = regexp.MustCompile(`Bad Magic Number`)
)

type errBadMagicNumber struct {
	count int
}

func (err errBadMagicNumber) Error() string {
	return fmt.Sprintf("bad magic number (%d counted)", err.count)
}

func waitForUBoot() {
	log.Println(term.Green("Waiting to start boot process"))
	expectOrPanic("==============================================", time.Minute)

	log.Println(term.Green("Waiting for UBoot"))
	expectOrPanic("UBoot Version", time.Second*10)
}

// verifyBoot verifies if the kernel boots successfully.
func verifyBoot() error {
	log.Println(term.Green("Waiting for booting image"))
	expectOrPanic("Booting image at", time.Second*5)
	log.Println(term.Green("Booting image"))

	_, _, _, err := e.ExpectSwitchCase([]expect.Caser{
		&expect.Case{
			// everything fine in this case
			R: regexp.MustCompile(`Please press Enter to activate this console`),
			T: expect.OK(),
		},
		&expect.Case{
			// not good
			R: reBadMagicNumber,
			T: expect.Continue(statusFlashCorrupt),
		},
	}, time.Second*5)

	// bad magic number found?
	if errors.Is(err, statusFlashCorrupt) {
		count := 1
		for {
			// Read all following bad magic numbers
			if _, _, err := e.Expect(reBadMagicNumber, time.Second); err == nil {
				count++
			} else {
				break
			}
		}

		return &errBadMagicNumber{count}
	}

	return err
}

func flashAndReboot() {
	log.Println(term.Green(`Choose operation "write to Flash via TFTP"`))

	_, err := e.ExpectBatch([]expect.Batcher{
		&expect.BExp{R: "Please choose the operation:"},
		&expect.BExp{R: "2: Load system code then write to Flash via TFTP"},
		&expect.BExp{R: "default: "},
		&expect.BSnd{S: "2"},
		&expect.BExp{R: "Are you sure\\?"},
		&expect.BSnd{S: "y"},
		&expect.BExp{R: "Input device IP.* ==:"},
		&expect.BSnd{S: *deviceIP + "\n"},
		&expect.BExp{R: "Input server IP.* ==:"},
		&expect.BSnd{S: *localIP + "\n"},
		&expect.BExp{R: "Input Linux Kernel filename .* ==:"},
		&expect.BSnd{S: kernelFilename + "\n"},
		&expect.BExp{R: "ETH_STATE_ACTIVE"},
		&expect.BExp{R: "Loading:"},
	}, time.Second*30)
	if err != nil {
		panic(err)
	}

	log.Println(term.Green("Waiting for TFTP connection"))

	_, err = e.ExpectBatch([]expect.Batcher{
		&expect.BCas{C: []expect.Caser{
			&expect.Case{
				R: regexp.MustCompile(`Got it`),
				T: expect.OK(),
			},
			&expect.Case{
				R:  regexp.MustCompile(`Retry count exceeded; starting again`),
				T:  expect.Continue(expect.NewStatus(0, "retry count exceeded")),
				Rt: 10,
			},
		}},
	}, time.Minute*10)
	if err != nil {
		panic(err)
	}

	expectOrPanic(fmt.Sprintf("Bytes transferred = %d", filesize), time.Minute)

	log.Println(term.Green("Flashing"))
	expectOrPanic("Done!", time.Minute)

	log.Println(term.Green("Rebooting"))

	_, err = e.ExpectBatch([]expect.Batcher{
		&expect.BExp{R: "Starting kernel \\.\\.\\."},
		&expect.BExp{R: "Please press Enter to activate this console"},
	}, time.Minute)
	if err != nil {
		panic(err)
	}

	log.Println(term.Green("Activating console"))
	e.Send("\n")
	expectOrPanic("root@(.+):/#", time.Second)

	e.Send("reboot\n")
}

func expectOrPanic(pattern string, timeout time.Duration) {
	_, _, err := e.Expect(regexp.MustCompile(pattern), timeout)
	if err != nil {
		panic(err)
	}
}
