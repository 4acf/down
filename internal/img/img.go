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

	/* example pixel-by-pixel read
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r >>= 8
			g >>= 8
			b >>= 8
		}
	}
	*/

	return img, nil

}
