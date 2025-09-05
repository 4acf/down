package main

import (
	"down/internal/img"
	"down/internal/soundfile"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {

	inputFilepath := flag.String("i", "", "input filepath (can be file or directory)")
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

	info, err := os.Stat(*inputFilepath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("file or directory %s does not exist", *inputFilepath)
		return
	}

	imgPaths, err := img.GetImgPaths(info.IsDir(), *inputFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, path := range imgPaths {
		img, err := img.Read(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		soundfile := soundfile.NewSoundfile(&img, strconv.Itoa(index))
		err = soundfile.Wav(*outputDirectory)
		if err != nil {
			fmt.Printf("wav file was not fully completed: %s", err)
		}
	}

	//spectrogram

}
