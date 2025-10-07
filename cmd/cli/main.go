package main

import (
	"down/config"
	"down/internal/imagefile"
	"down/internal/soundfile"
	"down/internal/spectrogram"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	AUDIO_OUTPUT_DIR       = "audio"
	SPECTROGRAM_OUTPUT_DIR = "spectrogram"
)

func main() {

	inputFilepath := flag.String("i", "", "input filepath (can be file or directory)")
	outputDirectory := flag.String("o", "./output", "output directory")
	progressEnabled := flag.Bool("p", true, "display progress, set using -p=true|false, default true")

	flag.Parse()

	if *inputFilepath == "" {
		fmt.Print("[no input filepath specified] input filepath is required, please specify input filepath with the -i flag")
		return
	}

	inputInfo, err := os.Stat(*inputFilepath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("[could not locate input directory] file or directory %s does not exist", *inputFilepath)
		return
	}

	audioOutputDirectory := filepath.Join(*outputDirectory, AUDIO_OUTPUT_DIR)
	spectrogramOutputDirectory := filepath.Join(*outputDirectory, SPECTROGRAM_OUTPUT_DIR)

	_, err = os.Stat(*outputDirectory)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(*outputDirectory, 0755)
		if err != nil {
			fmt.Printf("[failed to create output directory] %s", err)
			return
		}

		err = os.Mkdir(audioOutputDirectory, 0755)
		if err != nil {
			fmt.Printf("[failed to create output audio directory] %s", err)
			return
		}

		err = os.Mkdir(spectrogramOutputDirectory, 0755)
		if err != nil {
			fmt.Printf("[failed to create output audio directory] %s", err)
			return
		}
	}

	config := config.NewConfig(audioOutputDirectory, spectrogramOutputDirectory, *progressEnabled)

	//this function can return imagefiles with duplicate names, consider outputting files in directory tree(s) that matches the input
	reader := imagefile.NewReader(&config)
	imagefiles, err := reader.GetImagefiles(inputInfo, *inputFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, imagefile := range imagefiles {
		img, err := imagefile.Read()
		if err != nil {
			fmt.Printf("\n[%s read failed] %s", imagefile.Name(), err)
			continue
		}

		soundfile := soundfile.NewSoundfile(&config, &img, imagefile.Name())
		err = soundfile.Wav()
		if err != nil {
			fmt.Printf("\n[%s audio writeout failed] %s", soundfile.Name(), err)
			continue
		}

		spectrogram := spectrogram.NewSpectrogram(&config)
		colorlessImage, err := spectrogram.Image(soundfile.Data(), img.Bounds())
		if err != nil {
			fmt.Printf("\n[%s spectrogram calculation failed] %s", soundfile.Name(), err)
			continue
		}

		err = imagefile.Write(colorlessImage)
		if err != nil {
			fmt.Printf("\n[%s image writeout failed] %s", imagefile.Name(), err)
			continue
		}
	}

}
