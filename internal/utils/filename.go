package utils

import "path/filepath"

func RemoveExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}
