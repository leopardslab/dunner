package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetConfigs(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"

	var content = []byte(`test:
    - image: node
      command: ["node", "--version"]`)

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
		Name:    "",
		Image:   "node",
		Command: []string{"node", "--version"},
	}
	var tasks = make(map[string][]Task)
	tasks["test"] = []Task{task}
	var expected = Configs{
		Tasks: tasks,
	}

	imgs, _ := expected.GetAllImages()
	if !reflect.DeepEqual([]string{task.Image}, imgs) {
		t.Fatalf("Images list not equal to expected")
	}

	if !reflect.DeepEqual(expected, *pout) {
		t.Fatalf("Output not equal to expected")
	}

}
