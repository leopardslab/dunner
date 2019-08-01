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
