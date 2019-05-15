package config

import (
	"fmt"
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

	errs := configs.Validate()

	if len(errs) != 0 {
		t.Fatalf("Configs Validation failed, expected to pass. got: %s", errs)
	}
}

func TestConfigs_ValidateWithNoTasks(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 1 {
		t.Fatalf("Configs validation failed, expected 1 error, got %s", errs)
	}
	expected := "Tasks must contain at least 1 item"
	if errs[0].Error() != expected {
		t.Fatalf("expected: %s, got: %s", expected, errs[0].Error())
	}
}

func TestConfigs_ValidateWithEmptyImageAndCommand(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := Task{Image: "", Command: []string{""}}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d : %s", len(errs), errs)
	}

	expected1 := "task 'stats': image is a required field"
	expected2 := "task 'stats': command[0] is a required field"
	if errs[0].Error() != expected1 {
		t.Fatalf("expected: %s, got: %s", expected1, errs[0].Error())
	}
	if errs[1].Error() != expected2 {
		t.Fatalf("expected: %s, got: %s", expected2, errs[1].Error())
	}
}

func TestConfigs_ValidateWithInvalidMountFormat(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Mounts = []string{"invalid_dir"}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d : %s", len(errs), errs)
	}

	expected := "task 'stats': mount directory 'invalid_dir' is invalid. Check format is '<valid_src_dir>:<valid_dest_dir>:<mode>' and has right permission level"
	if errs[0].Error() != expected {
		t.Fatalf("expected: %s, got: %s", expected, errs[0].Error())
	}
}

func TestConfigs_ValidateWithValidMountDirectory(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	wd, _ := os.Getwd()
	task.Mounts = []string{fmt.Sprintf("%s:%s:w", wd, wd)}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if errs != nil {
		t.Fatalf("expected no errors, got %s", errs)
	}
}

func TestConfigs_ValidateWithNoModeGiven(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	wd, _ := os.Getwd()
	task.Mounts = []string{fmt.Sprintf("%s:%s", wd, wd)}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if errs != nil {
		t.Fatalf("expected no errors, got %s", errs)
	}
}

func TestConfigs_ValidateWithInvalidMode(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	wd, _ := os.Getwd()
	task.Mounts = []string{fmt.Sprintf("%s:%s:ab", wd, wd)}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	expected := fmt.Sprintf("task 'stats': mount directory '%s' is invalid. Check format is '<valid_src_dir>:<valid_dest_dir>:<mode>' and has right permission level", task.Mounts[0])
	if errs[0].Error() != expected {
		t.Fatalf("expected: %s, got: %s", expected, errs[0].Error())
	}
}

func TestConfigs_ValidateWithInvalidMountDirectory(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Mounts = []string{"blah:foo:w"}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d : %s", len(errs), errs)
	}

	expected := "task 'stats': mount directory 'blah:foo:w' is invalid. Check format is '<valid_src_dir>:<valid_dest_dir>:<mode>' and has right permission level"
	if errs[0].Error() != expected {
		t.Fatalf("expected: %s, got: %s", expected, errs[0].Error())
	}
}

func getSampleTask() Task {
	return Task{Image: "image_name", Command: []string{"node", "--version"}}
}

func TestInitValidatorForNilTranslation(t *testing.T) {
	vals := []customValidation{{tag: "foo", translation: "", validationFn: nil}}

	err := initValidator(vals)

	expected := "failed to register validation: Function cannot be empty"
	if err == nil {
		t.Fatalf("expected %s, got %s", expected, err)
	}
	if err.Error() != expected {
		t.Fatalf("expected %s, got %s", expected, err.Error())
	}
}
