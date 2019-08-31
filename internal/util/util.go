package util

import (
	"os"
	"os/user"
	"path"
	"strings"
)

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	HomeDir = user.HomeDir
}

// HomeDir is the environment variable HOME
var HomeDir string

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
