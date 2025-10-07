package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	FILE_PERMISSIONS = 0755
)

func InitializeOutputDirectories(outputDirectory, audioOutputDirectory, spectrogramOutputDirectory string) error {
	_, err := os.Stat(outputDirectory)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(outputDirectory, FILE_PERMISSIONS)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(audioOutputDirectory, FILE_PERMISSIONS)
	if err != nil {
		return err
	}

	err = os.MkdirAll(spectrogramOutputDirectory, FILE_PERMISSIONS)
	if err != nil {
		return err
	}
	return nil
}

func CreateFinalPath(parent, internal, filename string) (string, error) {
	paths := sanitizePaths(internal)

	err := createPaths(parent, paths)
	if err != nil {
		return "", err
	}

	return filepath.Join(parent, pathsToString(paths), filename), nil
}

func sanitizePaths(internal string) []string {
	n := len(internal)
	var result []string
	var sb strings.Builder

	for i := 0; i <= n; i++ {
		if i == n || internal[i] == os.PathSeparator {
			result = append(result, sb.String())
			sb.Reset()
			continue
		}
		sb.WriteByte(internal[i])
	}

	m := len(result)
	if m <= 2 {
		return []string{}
	}
	return result[1 : m-1]
}

func createPaths(parent string, paths []string) error {
	current := parent
	for _, path := range paths {
		current := filepath.Join(current, path)
		err := os.MkdirAll(current, FILE_PERMISSIONS)
		if err != nil {
			return err
		}
	}
	return nil
}

func pathsToString(paths []string) string {
	return filepath.Join(paths...)
}
