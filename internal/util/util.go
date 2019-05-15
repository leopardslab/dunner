package util

import (
	"os"
	"path"
	"strings"
)

// HomeDir is the environment variable HOME
var HomeDir = os.Getenv("HOME")

// DirExists returns true if the given param is a valid existing directory
func DirExists(dir string) bool {
	if strings.HasPrefix(dir, "~") {
		dir = path.Join(HomeDir, strings.Trim(dir, "~"))
	}
	src, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return src.IsDir()
}
