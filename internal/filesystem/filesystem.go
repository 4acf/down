package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

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
		err := os.MkdirAll(current, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func pathsToString(paths []string) string {
	return filepath.Join(paths...)
}
