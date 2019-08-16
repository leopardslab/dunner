package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/leopardslab/dunner/internal"
	"github.com/leopardslab/dunner/internal/util"
	"github.com/leopardslab/dunner/pkg/docker"
	"github.com/spf13/viper"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestGetConfigs(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"

	if err := os.Setenv("MYDUNNER", "dunner"); err != nil {
		t.Fatal(err)
	}
	defer os.Setenv("MYDUNNER", "")

	var content = []byte(`
test:
  - image: node:10.15.0
    user: 20
    commands:
      - ["node", "--version"]
      - ["npm", "--version"]
    envs:
      - MYVAR=MYVAL
      - MYUSR=` + "`$MYDUNNER`")

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
		Image:    "node:10.15.0",
		Commands: [][]string{{"node", "--version"}, {"npm", "--version"}},
		User:     "20",
		Envs:     []string{"MYVAR=MYVAL", "MYUSR=dunner"},
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

func TestParseEnv_InvalidEnv(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Image = "node:10.15.0"
	task.Envs = []string{"MYVAR=MYVAL", "MYUSR=dunner=invalid"}
	tasks["test"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	expectedErr := fmt.Errorf(
		`config: invalid format of environment variable: %s`,
		"MYUSR=dunner=invalid",
	)

	if err := ParseEnv(configs); err.Error() != expectedErr.Error() {
		t.Fatalf("Did not receive proper error on invalid format of environment variable, %v != %v", err, expectedErr)
	}
}

func TestParseEnv_EnvNotExist(t *testing.T) {
	tasks := make(map[string][]Task, 0)
	task := getSampleTask()
	task.Image = "node:10.15.0"
	task.Envs = []string{"MYVAR=MYVAL", "MYUSR=`$MYDUNNER`"}
	tasks["test"] = []Task{task}
	configs := &Configs{Tasks: tasks}

	expectedErr := fmt.Errorf(
		`config: could not find environment variable '%v' in %s file or among host environment variables`,
		"MYDUNNER",
		viper.GetString("DotenvFile"),
	)

	if err := ParseEnv(configs); err.Error() != expectedErr.Error() {
		t.Fatalf("Did not receive proper error on invalid format of environment variable, %v != %v", err, expectedErr)
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

func TestGetDunnerTaskFileWithCustomFileFromUser(t *testing.T) {
	taskFile := ".test_dunner.yaml"

	got, err := getDunnerTaskFile(taskFile)

	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if got != taskFile {
		t.Fatalf("expected original taskfile from user %s, got %s", taskFile, got)
	}
}

func TestGetDunnerTaskFileWithDefaultValue(t *testing.T) {
	taskFile := internal.DefaultDunnerTaskFileName

	got, err := getDunnerTaskFile(taskFile)

	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if !strings.HasSuffix(got, taskFile) {
		t.Fatalf("expected taskfile to end with %s, got %s", taskFile, got)
	}
}

func TestGetConfigsWhenNotPresentTillRoot(t *testing.T) {
	taskFile := internal.DefaultDunnerTaskFileName
	revert := setup(t)
	defer revert()

	got, err := GetConfigs(taskFile)

	if got != nil {
		t.Errorf("expected Configs to be nil, got %s", got)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectedErr := "failed to find Dunner task file"
	if err.Error() != expectedErr {
		t.Fatalf("expected error: %s, got: %s", expectedErr, err.Error())
	}
}

func setup(t *testing.T) func() {
	folder, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Failed to create temp dir: %s", err.Error())
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get working directory: %s", err.Error())
	}

	if err = os.Chdir(folder); err != nil {
		t.Errorf("Failed to change working directory: %s", err.Error())
	}
	return func() {
		if err = os.Chdir(previous); err != nil {
			t.Errorf("Failed to revert change in working directory: %s", err.Error())
		}
	}
}
