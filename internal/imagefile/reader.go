package imagefile

import (
	"down/config"
	"down/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type Reader struct {
	config *config.Config
}

func NewReader(config *config.Config) Reader {
	return Reader{
		config: config,
	}
}

func (reader *Reader) GetImagefiles(info os.FileInfo, path string) ([]Imagefile, error) {
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
				imagefiles = append(imagefiles, Imagefile{config: reader.config, path: p, name: name})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if isImageFile(path) {
			name := utils.RemoveExtension(info.Name())
			imagefiles = append(imagefiles, Imagefile{config: reader.config, path: path, name: name})
		}
	}

	return imagefiles, nil
}
