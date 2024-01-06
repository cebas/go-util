package util

import (
	"os"
	"path/filepath"
)

func ExecutableDir() (dir string) {
	// get executable path
	exePath, err := os.Executable()
	FatalErrorCheck(err)

	// get executable folder
	dir = filepath.Dir(exePath)

	return
}
