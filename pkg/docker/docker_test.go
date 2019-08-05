package docker

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/leopardslab/dunner/internal/settings"
)

func TestStep_Exec(t *testing.T) {
	var testNodeVersion = "10.15.0"
	results, err := runCommand([]string{"node", "--version"}, "./", testNodeVersion)
	if err != nil {
		t.Error(err)
	}

	strOut := (*results)[0].Output
	var res = strings.Trim(strings.Split(strOut, "v")[1], "\n")
	if res != testNodeVersion {
		t.Fatalf("Detected version of node container: '%s'; Expected output: '%s'", res, testNodeVersion)
	}
}

func TestStep_Exec_WorkingDirAbs(t *testing.T) {
	var testNodeVersion = "10.15.0"
	var absPath = "/go"
	results, err := runCommand([]string{"pwd"}, absPath, testNodeVersion)
	if err != nil {
		t.Error(err)
	}

	res := strings.Trim((*results)[0].Output, "\n")
	//var res = strings.Trim(strings.Split(strOut, "v")[1], "\n")
	if res != absPath {
		t.Fatalf("Detected working directory of node container: '%s'; Expected output: '%s'",
			res,
			absPath,
		)
	}
}

func TestStep_Exec_WorkingDirRel(t *testing.T) {
	var testNodeVersion = "10.15.0"
	var relPath = "./"
	results, err := runCommand([]string{"pwd"}, relPath, testNodeVersion)
	if err != nil {
		t.Error(err)
	}

	res := strings.Trim((*results)[0].Output, "\n")
	//var res = strings.Trim(strings.Split(strOut, "v")[1], "\n")
	if res != filepath.Join("/dunner", relPath) {
		t.Fatalf(
			"Detected working directory of node container: '%s'; Expected output: '%s'",
			res,
			filepath.Join("/dunner", relPath),
		)
	}
}

func runCommand(command []string, dir string, nodeVer string) (*[]Result, error) {
	settings.Init()
	step := &Step{
		Task:    "test",
		Name:    "node",
		Image:   "node:" + nodeVer,
		Command: command,
		Env:     nil,
		Volumes: nil,
		WorkDir: dir,
	}

	return step.Exec()
}
