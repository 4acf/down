package img

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func GetImgPaths(isDir bool, path string) ([]string, error) {
	var paths []string

	if isDir {
		return nil, errors.New("not implemented yet")
	} else {
		paths = append(paths, path)
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
