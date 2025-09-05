package img

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func GetImgPaths(isDir bool, path string) ([]string, error) {
	var paths []string

	isImageFile := func(p string) bool {
		fileExtension := strings.ToLower(filepath.Ext(p))
		return fileExtension == ".jpg" || fileExtension == ".png"
	}

	if isDir {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isImageFile(p) {
				paths = append(paths, p)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if isImageFile(path) {
			paths = append(paths, path)
		}
	}

	return paths, nil
}

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
