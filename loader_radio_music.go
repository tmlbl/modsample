package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ Loader = &RadioMusicLoader{}

// Files should be:
// Mono
// 16 Bit
// 44.1 kHz
// Headerless Wav files
// using the .raw suffix
//
// The card should be structured as 16 folders, named "0", "1", "2" etc
// (without the quote marks) Inside the folders, files can be named however you
// like. Firmwares before the 2017 version required filenames to be in 8:3
// format: NOISE.RAW, GOAT.RAW, HPSCHD.RAW. This is no longer the case with the
// 2017 firmware.
//
// MAXIMUM FILE LIMIT: You can add 48 files into each of the 16 folders.
// However, the module cannot handle more than about 330 files in total. This
// limit may be lifted with a future firmware upgrade.
type RadioMusicLoader struct {
	opts RadioMusicLoaderOptions
}

type RadioMusicLoaderOptions struct {
	RootPath string
}

func NewRadioMusicLoader(opts RadioMusicLoaderOptions) (*RadioMusicLoader, error) {
	l := &RadioMusicLoader{
		opts: opts,
	}
	err := l.ensureDirectoryStructure()
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Make sure we have the 16 folders named 0-15
func (l *RadioMusicLoader) ensureDirectoryStructure() error {
	for i := 0; i < 16; i++ {
		path := filepath.Join(l.opts.RootPath, fmt.Sprint(i))
		if _, err := os.Stat(path); err != nil {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Find the next usable sample location and create a path
func (l *RadioMusicLoader) getNextPathName() (string, error) {
	for i := 0; i < 16; i++ {
		path := filepath.Join(l.opts.RootPath, fmt.Sprint(i))
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return "", err
		}
		// If there is room in the bank, just name the file by combining
		// bank number and number of files in the bank
		if len(files) < 48 {
			return filepath.Join(path,
				fmt.Sprintf("%d%d.RAW", i, len(files))), nil
		}
	}
	// No room on the card!
	return "", ErrCardFull
}

func (l *RadioMusicLoader) AddSample(path string) (error, bool) {
	// ffmpeg -i 111.mp3 -acodec pcm_s16le -ac 1 -ar 16000 out.wav
	args := []string{"-i", path,
		// 16 bit PCM
		"-acodec", "pcm_s16le",
		// Converts it to mono
		"-ac", "1",
		// 44.1khz
		"-ar", "44100",
	}
	out, err := l.getNextPathName()
	if err != nil {
		if err == ErrCardFull {
			return nil, true
		}
		return err, false
	}
	args = append(args, "-out", out)
	err = runForeground("ffmpeg", args)
	if err != nil {
		return err, false
	}
	return nil, true
}
