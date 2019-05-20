package docker

import (
	"strings"
	"testing"

	"github.com/leopardslab/dunner/internal/settings"
)

func TestStep_Do(t *testing.T) {
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

	results, err := step.Exec()
	if err != nil {
		t.Error(err)
	}

	strOut := (*results)[0].Output
	var result = strings.Trim(strings.Split(strOut, "v")[1], "\n")
	if result != testNodeVersion {
		t.Fatalf("Detected version of node container: '%s'; Expected output: '%s'", result, testNodeVersion)
	}
}
