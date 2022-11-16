package handler

import (
	"fmt"
	"os"
	"path/filepath"
)

const confFile = "goj.yml"

func findRoot(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("%s not found", confFile)
	}
	if _, err := os.Stat(filepath.Join(dir, confFile)); err == nil {
		return dir, nil
	}

	if dir == "/" {
		return "", fmt.Errorf("%s not found", confFile)
	}
	return findRoot(filepath.Dir(dir))
}
