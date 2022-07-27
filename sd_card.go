package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

func detectCards() ([]string, error) {
	var options []string

	// For Linux list /dev and look for mmcblk devices
	if runtime.GOOS == "linux" {
		files, err := ioutil.ReadDir("/dev")
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "mmc") {
				// Just get the device, not the partitions
				if !strings.Contains(f.Name(), "p") {
					options = append(options, f.Name())
				}
			}
		}
	}

	return options, nil
}

func formatCard(device string) error {
	// Check if the disk is mounted
	parts, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	for _, p := range parts {
		if strings.Contains(p.Device, device) {
			fmt.Println("Unmounting device", device,
				"from", p.Mountpoint)
			err = runForeground("sudo", []string{
				"umount", p.Mountpoint})
			if err != nil {
				return err
			}
		}
	}

	// sudo mkfs.vfat -I --mbr=y /dev/mmcblk0
	// -I forces deletion of existing partitions (and everything else)
	args := []string{
		"mkfs.fat", "-I", "--mbr=y",
		fmt.Sprintf("/dev/%s", device),
	}
	return runForeground("sudo", args)
}
