package main

import (
	"down/internal/imagefile"
	"down/internal/soundfile"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {

	inputFilepath := flag.String("i", "", "input filepath (can be file or directory)")
	outputDirectory := flag.String("o", "./", "output directory")
	progressEnabled := flag.Bool("p", true, "display progress, set using -p=true|false")

	/*
		debug := flag.Bool("d", false, "debug mode")
		audio := flag.Bool("a", true, "keep audio files")
		image := flag.Bool("j", true, "keep image files")
	*/

	flag.Parse()

	if *inputFilepath == "" {
		fmt.Print("input filepath is required, please specify input filepath with the -i flag")
		return
	}

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
		soundfile := soundfile.NewSoundfile(&img, imagefile.Name())
		err = soundfile.Wav(*outputDirectory, *progressEnabled)
		if err != nil {
			fmt.Printf("wav file was not fully completed: %s", err)
		}
	}

	//spectrogram

}
