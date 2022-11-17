package handler

import (
	"fmt"
	"os"
	"path/filepath"
)

const testDir = "test"
const confFile = "goj.yml"

func findRoot(curDir string) (string, error) {
	curDir = filepath.Clean(curDir)
	if _, err := os.Stat(curDir); err != nil {
		return "", fmt.Errorf("%s not found", confFile)
	}
	if _, err := os.Stat(filepath.Join(curDir, confFile)); err == nil {
		return curDir, nil
	}

	return findRoot(filepath.Join(curDir, ".."))
}
