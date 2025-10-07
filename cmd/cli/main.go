package main

import (
	"down/config"
	"down/internal/filesystem"
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

	audioOutputDirectory, spectrogramOutputDirectory, err := filesystem.InitializeOutputDirectories(*outputDirectory)
	if err != nil {
		fmt.Printf("[could not initialize output directories] %s", err)
		return
	}

	config := config.NewConfig(audioOutputDirectory, spectrogramOutputDirectory, *progressEnabled)

	reader := imagefile.NewReader(&config)
	imagefiles, err := reader.GetImagefiles(inputInfo, *inputFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, imagefile := range imagefiles {

		extensionlessName := imagefile.Name()

		imageContents, err := imagefile.Read()
		if err != nil {
			fmt.Printf("\n[%s read failed] %s", extensionlessName, err)
			continue
		}

		soundfile := soundfile.NewSoundfile(&config, &imageContents, &imagefile)
		err = soundfile.Wav()
		if err != nil {
			fmt.Printf("\n[%s wav writeout failed] %s", extensionlessName, err)
			continue
		}

		spectrogram := spectrogram.NewSpectrogram(&config)
		colorlessImage, err := spectrogram.Image(soundfile.Data(), imageContents.Bounds())
		if err != nil {
			fmt.Printf("\n[%s spectrogram calculation failed] %s", extensionlessName, err)
			continue
		}

		err = imagefile.Write(colorlessImage)
		if err != nil {
			fmt.Printf("\n[%s image writeout failed] %s", extensionlessName, err)
			continue
		}
	}

}
