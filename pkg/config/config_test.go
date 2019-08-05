package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/leopardslab/dunner/internal/util"
	"github.com/leopardslab/dunner/pkg/docker"
	"gopkg.in/go-playground/validator.v9"
)

func TestGetConfigs(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"

	var content = []byte(`
test:
  - image: node
    user: 20
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
		User:     "20",
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

	if len(errs) != 0 {
		t.Fatalf("Configs validation failed, expected no error, got %s", errs)
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

	expected1 := "task 'stats': image is required, unless the task has a `follow` field"
	expected2 := "task 'stats': command[0] is a required field"
	if errs[0].Error() != expected1 {
		t.Fatalf("expected: %s, got: %s", expected1, errs[0].Error())
	}
	if errs[1].Error() != expected2 {
		t.Fatalf("expected: %s, got: %s", expected2, errs[1].Error())
	}
}

func TestConfigs_ValidateForAliasTask(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	tasks["foo"] = []Task{Task{Image: "golang", Command: []string{"go", "version"}}}
	tasks["stats"] = []Task{Task{Follow: "foo"}}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %d : %s", len(errs), errs)
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

	expected := "task 'stats': mount directory 'invalid_dir' is invalid. Check format is '<valid_src_dir>:<valid_dest_dir>:<optional_mode>' and has right permission level"
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

func TestConfigs_ValidateWithMountDirFromEnv(t *testing.T) {
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

	expected := fmt.Sprintf("task 'stats': mount directory '%s' is invalid. Check format is '<valid_src_dir>:<valid_dest_dir>:<optional_mode>' and has right permission level", task.Mounts[0])
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

	expected := "task 'stats': mount directory 'blah:foo:w' is invalid. Check if source directory path exists."
	if errs[0].Error() != expected {
		t.Fatalf("expected: %s, got: %s", expected, errs[0].Error())
	}
}

func TestConfigs_ValidateWithValidEnvInMountDir(t *testing.T) {
	os.Setenv("TEST_DIR", util.HomeDir)
	defer os.Setenv("TEST_DIR", "")
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Mounts = []string{"`$TEST_DIR`:foo:w"}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %d : %s", len(errs), errs)
	}
}

func TestConfigs_ValidateWithEnvInMountDir_Invalid(t *testing.T) {
	os.Setenv("TEST_DIR", "/test_invalid")
	defer os.Setenv("TEST_DIR", "")
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Mounts = []string{"`$TEST_DIR`:foo:w"}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d : %s", len(errs), errs)
	}

	expected := "task 'stats': mount directory '`$TEST_DIR`:foo:w' is invalid. Check if source directory path exists."
	if errs[0].Error() != expected {
		t.Fatalf("expected: %s, got: %s", expected, errs[0].Error())
	}
}

func TestConfigs_ValidateWithNonExistingEnvInMountDir(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Mounts = []string{"`$TEST_DIR_DUNNER`:foo:w"}
	tasks["stats"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	errs := configs.Validate()

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d : %s", len(errs), errs)
	}

	expected := "task 'stats': mount directory '`$TEST_DIR_DUNNER`:foo:w' is invalid. Check if source directory path exists."
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

	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}
}

func TestInitValidatorForEmptyTag(t *testing.T) {
	vals := []customValidation{{tag: "", translation: "",
		validationFn: func(context.Context, validator.FieldLevel) bool { return false }}}

	err := initValidator(vals)

	expected := "failed to register validation: Function Key cannot be empty"
	if err == nil {
		t.Fatalf("expected %s, got %s", expected, err)
	}
	if err.Error() != expected {
		t.Fatalf("expected %s, got %s", expected, err.Error())
	}
}

var lookupEnvtests = []struct {
	in  string
	out string
	err error
}{
	{"", "", nil},
	{"foo", "foo", nil},
	{"/foo/bar", "/foo/bar", nil},
	{"/foo/`$bar", "/foo/`$bar", nil},
	{util.HomeDir, util.HomeDir, nil},
	{"`$HOME`", util.HomeDir, nil},
	{"`$HOME`/foo", util.HomeDir + "/foo", nil},
	{"`$HOME`/foo/`$HOME`", util.HomeDir + "/foo/" + util.HomeDir, nil},
	{"`$INVALID_TEST`/foo", "`$INVALID_TEST`/foo", fmt.Errorf("could not find environment variable 'INVALID_TEST'")},
}

func TestLookUpDirectory(t *testing.T) {
	for _, tt := range lookupEnvtests {
		t.Run(tt.in, func(t *testing.T) {
			parsedDir, err := lookupDirectory(tt.in)
			if parsedDir != tt.out {
				t.Errorf("got %q, want %q", parsedDir, tt.out)
			}
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("got %q, want %q", err, tt.err)
			}
		})
	}
}

func TestDecodeMount(t *testing.T) {
	step := &docker.Step{}
	mounts := []string{fmt.Sprintf("%s:/app:r", util.HomeDir)}

	err := DecodeMount(mounts, step)

	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}
	if (*step).ExtMounts == nil {
		t.Fatalf("expected ExtMounts to be set, got nil")
	}
	if len((*step).ExtMounts) != 1 {
		t.Fatalf("expected ExtMounts to be of length 1, got %d", len((*step).ExtMounts))
	}
	if (*step).ExtMounts[0].Source != util.HomeDir {
		t.Fatalf("expected ExtMounts to be %s, got %s", util.HomeDir, (*step).ExtMounts[0].Source)
	}
}

func TestDecodeMountWithEnvironmentVariable(t *testing.T) {
	step := &docker.Step{}
	mounts := []string{"`$HOME`:`$HOME`"}

	err := DecodeMount(mounts, step)

	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}
	if (*step).ExtMounts == nil {
		t.Fatalf("expected ExtMounts to be set, got nil")
	}
	if len((*step).ExtMounts) != 1 {
		t.Fatalf("expected ExtMounts to be of length 1, got %d", len((*step).ExtMounts))
	}
	if (*step).ExtMounts[0].Source != util.HomeDir {
		t.Fatalf("expected ExtMounts Source to be %s, got %s", util.HomeDir, (*step).ExtMounts[0].Source)
	}
	if (*step).ExtMounts[0].Target != util.HomeDir {
		t.Fatalf("expected ExtMounts Source to be %s, got %s", util.HomeDir, (*step).ExtMounts[0].Target)
	}
}
