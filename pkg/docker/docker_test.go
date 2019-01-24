package docker

import (
	"bytes"
	"strings"
	"testing"
)

func TestStep_Do(t *testing.T) {

	var testNodeVersion = "10.15.0"

	step := &Step{
		Task:    "test",
		Name:    "node",
		Image:   "node:" + testNodeVersion,
		Command: []string{"node", "--version"},
		Env:     nil,
		WorkDir: "test",
		Volumes: nil,
	}
	images := []string{step.Image}
	if err := PullImages(&images); err != nil {
		t.Error(err)
	}

	pout, err := step.Do()
	if err != nil {
		t.Error(err)
	}
	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(*pout)
	if err != nil {
		t.Error(err)
	}

	strOut := buffer.String()
	var result = strings.Trim(strings.Split(strOut, "v")[1], "\n")
	if result != testNodeVersion {
		t.Fatalf("Detected version of node container: '%s'; Expected output: '%s'", result, testNodeVersion)
	}
}
