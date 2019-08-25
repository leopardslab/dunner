package dunner

import (
	"fmt"
	"io/ioutil"
	"os"
	os_user "os/user"
	"testing"

	"github.com/leopardslab/dunner/pkg/config"
	"github.com/leopardslab/dunner/pkg/docker"
	"github.com/spf13/viper"
)

func TestDo(t *testing.T) {

	var content = []byte(`
envs:
  - GLB=VARBL
tasks:
  test:
    envs:
      - GLB=VARBL2
      - MYVAR=GLBVAL
    steps:
      - image: busybox
        user: 20
        command: ["ls", "$1"]
        envs:
          - MYVAR=MYVAL`)

	if err := doContent(&content); err != nil {
		t.Fatal(err)
	}
}

func TestDo_VerboseAsync(t *testing.T) {
	async := viper.GetBool("Async")
	viper.Set("Async", true)
	verbose := viper.GetBool("Verbose")
	viper.Set("Verbose", true)

	defer viper.Set("Async", async)
	defer viper.Set("Verbose", verbose)

	TestDo(t)
}

func TestDo_WithFollow(t *testing.T) {

	var content = []byte(`
envs:
  - GLB=VARBL
tasks:
  test:
    envs:
      - GLB=VARBL2
      - MYVAR=GLBVAL
    steps:
      - image: busybox
        user: 20
        commands:
          - ["ls", "$1"]
        envs:
          - MYVAR=MYVAL
      - follow: test2
  test2:
    steps:
      - image: busybox
        command: ["pwd"]`)

	if err := doContent(&content); err != nil {
		t.Fatal(err)
	}
}

func doContent(content *[]byte) error {
	var tmpFilename = ".testdunner.yaml"

	tmpFile, err := ioutil.TempFile("", tmpFilename)
	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(*content); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	defaultTaskFile := viper.GetString("DunnerTaskFile")
	viper.Set("DunnerTaskFile", tmpFile.Name())
	defer viper.Set("DunnerTaskFile", defaultTaskFile)

	Do(nil, []string{"test", "/"})
	return nil
}

func TestExecTask(t *testing.T) {
	var step = config.Step{
		Name:     "",
		Image:    "busybox",
		Commands: [][]string{{"ls", "/"}, {"ls", "$1"}},
		Envs:     []string{"MYVAR=MYVAL"},
	}
	var tasks = make(map[string]config.Task)
	tasks["test"] = config.Task{Steps: []config.Step{step}}
	var configs = config.Configs{
		Tasks: tasks,
	}

	if err := ExecTask(&configs, "test", []string{"/dunner"}); err != nil {
		t.Fatal(err)
	}
}

func TestExecTaskWithParseError(t *testing.T) {
	step := config.Step{
		Image: "busybox",
		Dir:   "`$INVALID_USER_NONEXISTING`",
	}
	tasks := make(map[string]config.Task)
	tasks["test"] = config.Task{Steps: []config.Step{step}}
	configs := config.Configs{Tasks: tasks}

	err := ExecTask(&configs, "test", []string{})

	expectedErr := "could not find environment variable 'INVALID_USER_NONEXISTING'"
	if err == nil || err.Error() != expectedErr {
		t.Fatalf("expected error: %s, got %s", expectedErr, err)
	}
}

func TestExecTaskAsync(t *testing.T) {
	async := viper.GetBool("Async")
	viper.Set("Async", true)
	defer viper.Set("Async", async)

	TestExecTask(t)
}

func TestGetDunnerUserFromStep(t *testing.T) {
	expected := "test_user"
	step := config.Step{User: expected}

	user := getDunnerUser(step)

	if user != expected {
		t.Errorf("got: %s, want: %s", user, expected)
	}
}

func TestGetDunnerUserFromUserEnv(t *testing.T) {
	user, _ := os_user.Current()
	want := user.Uid

	got := getDunnerUser(config.Step{})

	if got != want {
		t.Errorf("got: %s, want: %s", user, want)
	}
}

func TestPassArgs_MultipleCommands(t *testing.T) {
	step := docker.Step{
		Commands: [][]string{{"ls", "$1"}, {"ls", "$2"}},
	}
	args := []string{"/"}
	err := PassArgs(&step, &args)
	expectedErr := fmt.Errorf(`dunner: insufficient number of arguments passed`)
	if err.Error() != expectedErr.Error() {
		t.Fatal("Improper or no error for insufficient number of arguments")
	}
}

func TestPassArgs_SingleCommand(t *testing.T) {
	step := docker.Step{
		Command: []string{"cp", "$1", "$2"},
	}
	args := []string{"/"}
	err := PassArgs(&step, &args)
	expectedErr := fmt.Errorf(`dunner: insufficient number of arguments passed`)
	if err.Error() != expectedErr.Error() {
		t.Fatal("Improper or no error for insufficient number of arguments")
	}
}
