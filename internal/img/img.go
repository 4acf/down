package img

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func Read(filepath string) (image.Image, error) {
	reader, err := os.Open(filepath)
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
