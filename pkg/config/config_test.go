package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetConfigs(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"

	var content = []byte(`
test:
  - image: node
    commands:
      - ["node", "--version"]
      - ["npm", "--version"]
    envs:
      - MYVAR=MYVAL`)

	tmpFile, err := ioutil.TempFile("", tmpFilename)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	pout, err := GetConfigs(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	var task = Task{
		Name:     "",
		Image:    "node",
		Commands: [][]string{{"node", "--version"}, {"npm", "--version"}},
		Envs:     []string{"MYVAR=MYVAL"},
	}
	var tasks = make(map[string][]Task)
	tasks["test"] = []Task{task}
	var expected = Configs{
		Tasks: tasks,
	}

	if !reflect.DeepEqual(expected, *pout) {
		t.Fatalf("Output not equal to expected; %v != %v", expected, *pout)
	}

}
