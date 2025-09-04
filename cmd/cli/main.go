package main

import (
	"down/internal/img"
	"down/internal/soundfile"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {

	inputFilepath := flag.String("i", "", "input filepath")
	outputDirectory := flag.String("o", "./", "output directory")

	/*
		debug := flag.Bool("d", false, "debug mode")
		progress := flag.Bool("p", true, "display progress")
		audio := flag.Bool("a", true, "keep audio files")
		image := flag.Bool("j", true, "keep image files")
	*/

	flag.Parse()

	if *inputFilepath == "" {
		fmt.Print("input filepath is required, please specify input filepath with the -i flag")
		return
	}
	if _, err := os.Stat(*inputFilepath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("file %s does not exist", *inputFilepath)
		return
	}

	//get rgb image data
	img, err := img.Read(*inputFilepath)
	if err != nil {
		//for now we print the error and return, in the future when this block is in a loop this will break instead of return
		fmt.Println(err)
		return
	}

	//writeout to a wav file
	soundfile := soundfile.NewSoundfile(&img, "test.wav")
	err = soundfile.Wav(*outputDirectory)
	if err != nil {
		fmt.Printf("wav file was not fully completed: %s", err)
		return
	}

	//spectrogram

}
