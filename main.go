package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

// func main() {
// 	folders := pflag.StringArrayP("files", "f", []string{},
// 		"Directories to load samples from")
// 	format := pflag.Bool("format", false, "Whether to format the card")

// 	pflag.Parse()

// 	fmt.Println("Selected folders", folders)

// 	options, err := detectCards()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	device, err := choose("Choose a device", options)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if *format {
// 		fmt.Println("Formatting", device)
// 		err = formatCard(device)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	pathList, err := buildSampleList(*folders)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(pathList)

// 	var devicePath string

// 	// Check if the disk is mounted
// 	parts, err := disk.Partitions(false)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, p := range parts {
// 		if strings.Contains(p.Device, device) {
// 			devicePath = p.Mountpoint
// 		}
// 	}

// 	for i, path := range pathList {
// 		args := []string{"-i", path, "-ac", "1",
// 			filepath.Join(devicePath, fmt.Sprintf("Samp%d.wav", i))}
// 		err = runForeground("ffmpeg", args)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

func main() {
	folders := pflag.StringArrayP("files", "f", []string{},
		"Directories to load samples from")
	dest := pflag.StringP("dest", "d", "",
		"Path to the mounted SD card")
	pflag.Parse()

	fmt.Println(*folders, *dest)
}

func choose(prompt string, options []string) (string, error) {
	fmt.Println(prompt)
	for i, o := range options {
		fmt.Printf("%d: %s\n", i, o)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	i, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return "", err
	}
	if (i - 1) > len(options) {
		return "", fmt.Errorf("index out of range")
	}
	return options[i], nil
}

func buildSampleList(folders []string) ([]string, error) {
	pathList := []string{}
	for _, folder := range folders {
		filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			pathList = append(pathList, path)
			return nil
		})
	}
	return pathList, nil
}
