package imagefile

import (
	"down/config"
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

	for i := range xRes {
		for j := range yRes {
			channelValue := uint8(data[i][j] * 0xff)
			color := color.NRGBA{
				R: channelValue, G: channelValue, B: channelValue, A: 0xff,
			}
			newImage.SetNRGBA(i, yRes-j, color)
		}
	}

	file, err := os.Create("spectrogram.png")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, newImage); err != nil {
		return err
	}

	return nil
}
