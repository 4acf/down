package utils

import (
	"path/filepath"
	"strings"
)

const (
	PNG_EXTENSION = ".png"
	WAV_EXTENSION = ".wav"
)

func RemoveExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func AppendPng(name string) string {
	return appendExtension(name, PNG_EXTENSION)
}

func AppendWav(name string) string {
	return appendExtension(name, WAV_EXTENSION)
}

func appendExtension(name string, extension string) string {
	name = strings.TrimSpace(name)

	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}

	if strings.HasSuffix(name, ".") {
		return name + extension
	} else if !strings.HasSuffix(strings.ToLower(name), extension) {
		return name + extension
	}
	return name
}
