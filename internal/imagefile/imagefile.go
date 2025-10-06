package imagefile

import (
	"down/internal/utils"
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Imagefile struct {
	path string
	name string
	img  *image.Image
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

	imagefile.img = &img

	return img, nil
}

func GetImagefiles(info os.FileInfo, path string) ([]Imagefile, error) {
	var imagefiles []Imagefile

	isImageFile := func(filename string) bool {
		fileExtension := strings.ToLower(filepath.Ext(filename))
		return fileExtension == ".jpg" || fileExtension == ".png"
	}

	if info.IsDir() {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isImageFile(p) {
				name := utils.RemoveExtension(info.Name())
				imagefiles = append(imagefiles, Imagefile{path: p, name: name, img: nil})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if isImageFile(path) {
			name := utils.RemoveExtension(info.Name())
			imagefiles = append(imagefiles, Imagefile{path: path, name: name, img: nil})
		}
	}

	return imagefiles, nil
}

func (imagefile *Imagefile) Write(data [][]complex128) error {

	if imagefile.img == nil {
		return errors.New("cannot access image as imagefile has not been read")
	}

	bounds := (*imagefile.img).Bounds()
	xRes := bounds.Dx()
	yRes := bounds.Dy()

	newImage := image.NewNRGBA(
		image.Rect(0, 0, xRes, yRes),
	)

	frames := len(data)
	if frames == 0 {
		return errors.New("spectrogram has insufficient frame information")
	}

	frequencyBins := len(data[0])
	if frequencyBins == 0 {
		return errors.New("spectrogram has insufficient frequency information")
	}
	frequencyBins /= 2

	clamp := func(value, min, max float64) float64 {
		ternary := func(condition bool, a, b float64) float64 {
			if condition {
				return a
			}
			return b
		}
		return ternary(value > max, max, ternary(min > value, min, value))
	}

	scaleInt := func(value, maxIn, maxOut int) int {
		return int(float64(value) / float64(maxIn) * float64(maxOut))
	}

	minDb := float64(-120)

	for i := range xRes {
		for j := range yRes {
			frame := scaleInt(i, xRes, frames)
			frequencyBin := scaleInt(j, yRes, frequencyBins)
			abs := math.Abs(real(data[frame][frequencyBin]))
			db := 20 * math.Log10(abs+1e-10)
			db = clamp(db, minDb, 0)
			normalizedDb := (db + (minDb * -1)) / (minDb * -1)
			color := color.NRGBA{
				R: 0xff, G: 0xff, B: 0xff, A: uint8(normalizedDb * 255),
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
