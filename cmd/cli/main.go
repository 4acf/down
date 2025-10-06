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
)

func main() {

	inputFilepath := flag.String("i", "", "input filepath (can be file or directory)")
	outputDirectory := flag.String("o", "./", "output directory")
	progressEnabled := flag.Bool("p", true, "display progress, set using -p=true|false")
	keepAudio := flag.Bool("a", true, "keep audio files")
	keepSpectrogram := flag.Bool("s", true, "keep spectrogram files")

	flag.Parse()

	if *inputFilepath == "" {
		fmt.Print("input filepath is required, please specify input filepath with the -i flag")
		return
	}

	config := config.NewConfig(*outputDirectory, *progressEnabled, *keepAudio, *keepSpectrogram)

	info, err := os.Stat(*inputFilepath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("file or directory %s does not exist", *inputFilepath)
		return
	}

	//this function can return imagefiles with duplicate names, consider outputting files in directory tree(s) that matches the input
	imagefiles, err := imagefile.GetImagefiles(info, *inputFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, imagefile := range imagefiles {
		img, err := imagefile.Read()
		if err != nil {
			fmt.Println(err)
			continue
		}

		soundfile := soundfile.NewSoundfile(&config, &img, imagefile.Name())
		err = soundfile.Wav()
		if err != nil {
			fmt.Printf("\n%s writeout was not fully completed: %s", soundfile.Name(), err)
		}

		spectrogram := spectrogram.NewSpectrogram(&config)
		colorlessImage, err := spectrogram.Image(soundfile.Data(), img.Bounds())
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = imagefile.Write(colorlessImage)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
