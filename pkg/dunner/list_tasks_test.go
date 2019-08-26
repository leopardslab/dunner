package dunner

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func Test_ListTasksWhenConfigFileNotFound(t *testing.T) {
	viper.Set("DunnerTaskFile", "fileThatDoesnotExit.yaml")
	defer viper.Reset()

	err := ListTasks()

	expected := "open fileThatDoesnotExit.yaml: no such file or directory"
	if err == nil {
		t.Fatalf("got: %s, want: %s", err, expected)
	}
	if err.Error() != expected {
		t.Fatalf("got: %s, want: %s", err.Error(), expected)
	}
}

func Test_ListTasksSuccess(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"
	var content = []byte(`
envs:
  - GLB=VARBL
tasks:
  setup:
    steps:
      - image: node
        commands:
          - ["node", "--version"]
          - ["npm", "--version"]
        envs:
          - MYVAR=MYVAL
  build:
    steps:
      - image: node
        command: ["ls"]`)

	tmpFile := createDunnerTaskFile(t, content, tmpFilename)
	defer os.Remove(tmpFile.Name())
	defer viper.Reset()

	err := ListTasks()

	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func ExampleListTasks_successWithAllTasksAsBullets() {
	var tmpFilename = ".testdunner.yaml"
	var content = []byte(`
tasks:
  setup:
    steps:
      - image: node
        command: []
  build:
    steps:
      - image: node
        command: []`)

	tmpFile, err := ioutil.TempFile("", tmpFilename)
	if err != nil {
		panic(err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		panic(err)
	}

	if err := tmpFile.Close(); err != nil {
		panic(err)
	}

	viper.Set("DunnerTaskFile", tmpFile.Name())
	defer viper.Reset()
	defer os.Remove(tmpFile.Name())

	err = ListTasks()

	if err != nil {
		panic(err)
	}

	// Unordered output: Available Dunner tasks:
	// • setup
	// • build
	// Run `dunner do <task_name>` to run a dunner task.
}

func Test_ListTasksSuccessNoTasks(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"
	var content = []byte("")

	tmpFile := createDunnerTaskFile(t, content, tmpFilename)
	defer os.Remove(tmpFile.Name())
	defer viper.Reset()

	err := ListTasks()

	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func createDunnerTaskFile(t *testing.T, content []byte, tmpFilename string) *os.File {
	tmpFile, err := ioutil.TempFile("", tmpFilename)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	viper.Set("DunnerTaskFile", tmpFile.Name())
	return tmpFile
}
