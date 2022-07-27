package main

import (
	"fmt"
	"path/filepath"
)

var _ Loader = &TipTopLoader{}

// Samples for the TipTop Audio One eurorack sample player
// From the manual:
//
// ONE plays any mono WAV format with a bit depth of 16 or 24 bits and a sample
// rate from 8 to 96 kHz. The native rates supported by ONE are 44.1, 48, 88.2
// and 96 kHz and these can be correctly quantized to the Western 12 tone scale.
// Other rates will use 48k as their time base and be pitched up or down
// accordingly.
// The names can be any alphanumeric character (A-Z, 0-9) but are limited to 13
// characters in length and the 13 characters includes the .wav extension.
// ONE scans the microSD card for all compatible on power up and makes a table
// of the files. ONE can store a maximum of 256 files in the table.
// The minimum length for is 1024 samples. For the most seamless looping
// (like waveforms) should be a multiple of this number.
type TipTopLoader struct {
	opts    TipTopLoaderOptions
	nLoaded int
}

type TipTopLoaderOptions struct {
	RootPath string
}

func NewTipTopLoader(opts TipTopLoaderOptions) (*TipTopLoader, error) {
	return &TipTopLoader{
		opts: opts,
	}, nil
}

func (l *TipTopLoader) AddSample(path string) (error, bool) {
	// Convert to mono
	args := []string{"-i", path, "-ac", "1",
		filepath.Join(l.opts.RootPath,
			fmt.Sprintf("Samp%d.wav", l.nLoaded))}
	err := runForeground("ffmpeg", args)
	if err != nil {
		return err, false
	}
	l.nLoaded++
	return nil, false
}
