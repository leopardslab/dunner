package dunner

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func Test_ListTasksWhenConfigFileNotFound(t *testing.T) {
	viper.Set("DunnerTaskFile", "fileThatDoesnotExit.yaml")

	err := ListTasks()

	expected := "open fileThatDoesnotExit.yaml: no such file or directory"
	if err == nil {
		t.Fatalf("got: %s, want: %s", err, expected)
	}
	if err.Error() != expected {
		t.Fatalf("got: %s, want: %s", err.Error(), expected)
	}
}

func Test_ListTasksWhenValidationFails(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"
	var content = []byte(`
setup:
  - image: node
    commands:
      - ["node", "--version"]
      - ["npm", "--version"]
    envs:
      - MYVAR=MYVAL
build:
  - command: []`)

	tmpFile := createDunnerTaskFile(t, content, tmpFilename)
	defer os.Remove(tmpFile.Name())

	err := ListTasks()

	expected := "validation failed"
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
setup:
  - image: node
    commands:
      - ["node", "--version"]
      - ["npm", "--version"]
    envs:
      - MYVAR=MYVAL
build:
  - image: node
    command: []`)

	tmpFile := createDunnerTaskFile(t, content, tmpFilename)
	defer os.Remove(tmpFile.Name())

	err := ListTasks()

	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func Test_ListTasksSuccessNoTasks(t *testing.T) {
	var tmpFilename = ".testdunner.yaml"
	var content = []byte("")

	tmpFile := createDunnerTaskFile(t, content, tmpFilename)
	defer os.Remove(tmpFile.Name())

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
