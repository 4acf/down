package imagefile

import (
	"down/config"
	"down/internal/filesystem"
	"down/internal/spectrogram"
	"down/internal/utils"
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

type Imagefile struct {
	config *config.Config
	path   string
	name   string
}

func (imagefile *Imagefile) Name() string {
	return imagefile.name
}

func (imagefile *Imagefile) Path() string {
	return imagefile.path
}

func (imagefile *Imagefile) Read() (image.Image, error) {
	reader, err := os.Open(imagefile.path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (imagefile *Imagefile) Write(data [][]float64) error {

	xRes := len(data)
	if xRes == 0 {
		return errors.New("cannot color an image that has no x resolution information")
	}

	yRes := len(data[0])
	if yRes == 0 {
		return errors.New("cannot color an image that has no y resolution information")
	}

	newImage := image.NewNRGBA(
		image.Rect(0, 0, xRes, yRes),
	)

	paletteLength := len(spectrogram.PALETTE) / 3
	for i := range xRes {
		for j := range yRes {
			colorIndex := int(data[i][j] * float64(paletteLength-1))
			color := color.NRGBA{
				R: spectrogram.PALETTE[colorIndex*3],
				G: spectrogram.PALETTE[(colorIndex*3)+1],
				B: spectrogram.PALETTE[(colorIndex*3)+2],
				A: 0xff,
			}
			newImage.SetNRGBA(i, yRes-j, color)
		}
	}

	filename := utils.AppendPng(imagefile.name)
	outputFilepath, err := filesystem.CreateFinalPath(imagefile.config.SpectrogramOutputDirectory(), imagefile.path, filename)
	if err != nil {
		return err
	}

	file, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, newImage); err != nil {
		return err
	}

	return nil
}
