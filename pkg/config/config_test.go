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
    command: ["node", "--version"]
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
		Name:    "",
		Image:   "node",
		Command: []string{"node", "--version"},
		Envs:    []string{"MYVAR=MYVAL"},
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

func TestConfigs_Validate(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	tasks["stats"] = []Task{getSampleTask()}
	configs := &Configs{Tasks: tasks}

	errs, ok := configs.Validate()

	if !ok || len(errs) != 0 {
		t.Fatalf("Configs Validation failed, expected to pass")
	}
}

func TestConfigs_ValidateWithNoTasks(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	configs := &Configs{Tasks: tasks}

	errs, ok := configs.Validate()

	if !ok || len(errs) != 1 {
		t.Fatalf("Configs validation failed")
	}
	if errs[0].Error() != "dunner: No tasks defined" {
		t.Fatalf("Configs Validation error message not as expected")
	}
}

func TestConfigs_ValidateWithParseErrors(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := Task{Image: "", Command: []string{}}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs, ok := configs.Validate()

	if ok || len(errs) != 2 {
		t.Fatalf("Configs validation failed")
	}

	if errs[0].Error() != "dunner: [stats] Image repository name cannot be empty" || errs[1].Error() != "dunner: [stats] Commands not defined for task with image " {
		t.Fatalf("Configs Validation error message not as expected")
	}
}

func getSampleTask() Task {
	return Task{Image: "image_name", Command: []string{"node", "--version"}}
}
