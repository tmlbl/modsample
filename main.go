package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func selectRandom(in []SampleInfo, n int) []SampleInfo {
	indexes := []int{}
	for i := 0; i < n; i++ {
		indexes = append(indexes, rand.Intn(len(in)))
	}
	out := make([]SampleInfo, n)
	for i, ix := range indexes {
		out[i] = in[ix]
	}
	return out
}

func selectLargerThan(in []SampleInfo, than int) []SampleInfo {
	out := []SampleInfo{}
	for _, s := range in {
		if s.Info.Size() > int64(than*1000) {
			out = append(out, s)
		}
	}
	return out
}

func main() {
	folders := pflag.StringArrayP("files", "f", []string{},
		"Directories to load samples from")
	// dest := pflag.StringP("dest", "d", "",
	// 	"Path to the mounted SD card")
	random := pflag.BoolP("random", "r", false,
		"Whether to select random samples")
	largerThan := pflag.Int("larger-than", -1,
		"Files larger than the given size (in kb)")
	// smallerThan := pflag.Int("smaller-than", -1,
	// 	"Files smaller than the given size (in kb)")
	pflag.Parse()

	list, err := buildSampleList(*folders)
	if err != nil {
		fatal(err)
	}

	if *largerThan > 0 {
		list = selectLargerThan(list, *largerThan)
	}

	if *random {
		rand.Seed(time.Now().UnixMilli())
		list = selectRandom(list, 5)
	}

	for _, ln := range list {
		fmt.Println(ln.Path, ln.Info.Size())
	}
}

func fatal(err error) {
	log.Fatalln(err)
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

type SampleInfo struct {
	Path string
	Info os.FileInfo
}

func buildSampleList(folders []string) ([]SampleInfo, error) {
	list := []SampleInfo{}
	for _, folder := range folders {
		filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			stat, err := os.Stat(path)
			if err != nil {
				fatal(err)
			}
			list = append(list, SampleInfo{
				Path: path,
				Info: stat,
			})
			return nil
		})
	}
	return list, nil
}
