package imagefile

import (
	"down/internal/utils"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

type Imagefile struct {
	path string
	name string
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
				imagefiles = append(imagefiles, Imagefile{path: p, name: name})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if isImageFile(path) {
			name := utils.RemoveExtension(info.Name())
			imagefiles = append(imagefiles, Imagefile{path: path, name: name})
		}
	}

	return imagefiles, nil
}
