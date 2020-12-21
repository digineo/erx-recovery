# ER-X recovery tool

This tools helps to recover broken Ubiquiti EdgeRouter X devices using a serial console and TFTP.
It also detects damaged flash when the kernel is not able to boot.

It has been tested under Ubuntu 20.04.

## Requirements

* Linux
* Serial console attached to the EdgeRouter X.
* `eth0` of the EdgeRouter X connected to the local machine.

## Usage

```
Usage of erx-recovery:
  -device-ip string
    	device IP address for the TFTP client, needs to be in the same network as the TFTP server (default "172.16.3.212")
  -image string
    	path to linux kernel image (default "/path/to/kernel")
  -local-ip string
    	local IP address for the TFTP server (default "172.16.3.210")
  -log string
    	directory for log files (default "./logs")
  -tty string
    	path to serial console (default "/dev/ttyUSB0")
  -verbose
    	Increase verbosity (default true)
```

## Detection of damaged flash

If important parts of the flash memory (NAND) are damaged,
the kernel is unable to boot. In this case you see something like this in the log file:

```
Please choose the operation:
   1: Load system code to SDRAM via TFTP.
   2: Load system code then write to Flash via TFTP.
   3: Boot system code via Flash (default).
   4: Entr boot command line interface.
   7: Load Boot Loader code then write to Flash via Serial.
   9: Load Boot Loader code then write to Flash via TFTP.
   r: Start TFTP recovery.
default: 3
 0

3: System Boot system code via Flash.
## Booting image at c0040000 ...
Bad Magic Number,53F511F6
Search header in next block address 460000
Bad Magic Number,578CD49B
Search header in next block address 480000
Bad Magic Number,D718AC07
Search header in next block address 4a0000
Bad Magic Number,F074335F
Search header in next block address 4c0000
Bad Magic Number,E761A28A
Search header in next block address 4e0000
Bad Magic Number,A8CE778A
Search header in next block address 500000
... [many more]
```
