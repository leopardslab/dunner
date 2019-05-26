package initialize

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/leopardslab/dunner/internal"
	"github.com/leopardslab/dunner/pkg/config"
	yaml "gopkg.in/yaml.v2"
)

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

func TestInitializeSuccess(t *testing.T) {
	revert := setup(t)
	defer revert()
	var filename = ".test_dunner.yml"
	if err := initProject(filename); err != nil {
		t.Errorf("Failed to open dunner task file %s: %s", filename, err.Error())
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Errorf("Failed to open dunner task file %s: %s", filename, err.Error())
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("Failed to read dunner task file %s: %s", filename, err.Error())
	}

	var configs config.Configs
	if err := yaml.Unmarshal(fileContents, &configs.Tasks); err != nil {
		t.Errorf("Task file config structure invalid: %s", err.Error())
	}
}

func TestInitializeWhenFileExists(t *testing.T) {
	revert := setup(t)
	defer revert()
	var filename = ".test_dunner.yml"
	createFile(t, filename, internal.DefaultTaskFileContents)

	expected := fmt.Sprintf("%s already exists", filename)
	err := initProject(filename)
	if err == nil {
		t.Errorf("expected: %s, got nil", expected)
	}
	if expected != err.Error() {
		t.Errorf("expected: %s, got: %s", expected, err.Error())
	}
}

func createFile(t *testing.T, filename, contents string) {
	if err := ioutil.WriteFile(filename, []byte(contents), 0644); err != nil {
		t.Errorf("Failed to create file: %s", err.Error())
	}
}
