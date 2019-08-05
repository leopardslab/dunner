package docker

import (
	"github.com/leopardslab/dunner/internal/settings"
)

func ExampleStep_Exec() {
	settings.Init()
	var testNodeVersion = "10.15.0"
	step := &Step{
		Task:     "test",
		Name:     "node",
		Image:    "node:" + testNodeVersion,
		Commands: [][]string{{"node", "--version"}},
		Env:      nil,
		Volumes:  nil,
	}

	err := step.Exec()
	if err != nil {
		panic(err)
	}
	// Output: OUT: v10.15.0
}

func ExampleStep_Exec_WorkingDirAbs() {
	var testNodeVersion = "10.15.0"
	var absPath = "/go"
	err := runCommand([]string{"pwd"}, absPath, testNodeVersion)

	if err != nil {
		panic(err)
	}
	// Output: OUT: /go
}

func ExampleStep_Exec_WorkingDirRel() {
	var testNodeVersion = "10.15.0"
	var relPath = "./"
	err := runCommand([]string{"pwd"}, relPath, testNodeVersion)
	if err != nil {
		panic(err)
	}
	// Output: OUT: /dunner
}

func runCommand(command []string, dir string, nodeVer string) error {
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