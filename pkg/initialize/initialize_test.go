package initialize

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/leopardslab/dunner/internal"
	"github.com/leopardslab/dunner/pkg/config"
	"github.com/leopardslab/dunner/pkg/global"

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

func TestInitProjectSuccess(t *testing.T) {
	revert := setup(t)
	defer revert()
	var filename = ".test_dunner.yml"
	if err := InitProject(filename, nil); err != nil {
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
	err := InitProject(filename, nil)
	if err == nil {
		t.Errorf("expected: %s, got nil", expected)
	}
	if expected != err.Error() {
		t.Errorf("expected: %s, got: %s", expected, err.Error())
	}
}

func TestInitializeFilenameIsInvalid(t *testing.T) {
	revert := setup(t)
	defer revert()
	var filename = "#Q$EJL_doesntexist/.test_dunner.yml"

	expected := fmt.Sprintf("open %s: no such file or directory", filename)
	err := InitProject(filename, nil)
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

func TestGetRecipeMetadataWhenMetadataDoesNotExist(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	_, err := getRecipeMetadata(server.URL)

	expected := fmt.Sprintf("Failed to download metadata: Error downloading file %s: 404 Not Found", server.URL)
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error: %s, got %s", expected, err)
	}
}

func TestGetRecipeMetadataWhenMetadataIsInvalidYaml(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "name=foo", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	metadata, err := getRecipeMetadata(server.URL)

	expected := "Failed to unmarshal metadata"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err)
	}
	if metadata != nil {
		t.Errorf("Expected metadata to be nil, got %s", metadata)
	}
}

func TestGetRecipeMetadataSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "name: foo\ndescription: foo_description\nversion: 0.0.1\npreInstallCmd: \"echo hello\"\npostInstallMessage: \"Done\"", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	metadata, err := getRecipeMetadata(server.URL)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if metadata == nil {
		t.Errorf("Expected metadata, got %s", metadata)
	}
	if (*metadata).Name != "foo" {
		t.Errorf("Expected foo, got %s", (*metadata).Name)
	}
	if (*metadata).Description != "foo_description" {
		t.Errorf("Expected foo_description, got %s", (*metadata).Description)
	}
	if (*metadata).Version != "0.0.1" {
		t.Errorf("Expected 0.0.1, got %s", (*metadata).Version)
	}
	if (*metadata).PreInstallCmd != "echo hello" {
		t.Errorf("Expected echo hello, got %s", (*metadata).PreInstallCmd)
	}
	if (*metadata).PostInstallMessage != "Done" {
		t.Errorf("Expected Done, got %s", (*metadata).PostInstallMessage)
	}
}

func TestConstructURLs(t *testing.T) {
	expectedMetadataURL := fmt.Sprintf("%sfoo/metadata.yml", global.DunnerCookbookBaseURL)
	got := getMetadataURL("foo")
	if got != expectedMetadataURL {
		t.Errorf("expected URL %s, got %s", expectedMetadataURL, got)
	}

	expectedDunnerTaskURL := fmt.Sprintf("%sfoo/.dunner.yaml", global.DunnerCookbookBaseURL)
	got = getDunnerTaskURLOfRecipe("foo")
	if got != expectedDunnerTaskURL {
		t.Errorf("expected URL %s, got %s", expectedDunnerTaskURL, got)
	}
}

func TestInitProjectWithRecipe(t *testing.T) {
	revert := setup(t)
	defer revert()
	wd, _ := os.Getwd()

	taskFilePath := fmt.Sprintf("%s/.test_init_dunner.yaml", wd)
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "name: foo", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	getMetadataURL = func(name string) string { return server.URL }
	getDunnerTaskURLOfRecipe = func(string) string { return server.URL }
	defer server.Close()

	err := InitProject(".test_init_dunner.yaml", []string{"foo"})

	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
	if _, err := os.Stat(taskFilePath); err != nil {
		t.Errorf("Dunner task file does not exist")
	}
}

func TestInitProjectWithRecipeWithPreInstallCmd(t *testing.T) {
	revert := setup(t)
	defer revert()
	wd, _ := os.Getwd()

	taskFilePath := fmt.Sprintf("%s/.test_init_dunner.yaml", wd)
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "name: foo\npreInstallCmd: \"echo hello\"\npostInstallMessage: \"Done\"", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	getMetadataURL = func(name string) string { return server.URL }
	getDunnerTaskURLOfRecipe = func(string) string { return server.URL }
	defer server.Close()

	err := InitWithRecipe(".test_init_dunner.yaml", "foo")

	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
	if _, err := os.Stat(taskFilePath); err != nil {
		t.Errorf("Dunner task file does not exist")
	}
}

func TestInitProjectWithRecipeWithInvalidPreInstallCmd(t *testing.T) {
	revert := setup(t)
	defer revert()
	wd, _ := os.Getwd()

	taskFilePath := fmt.Sprintf("%s/.test_init_dunner.yaml", wd)
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "name: foo\npreInstallCmd: \"invalid_cmd hello\"", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	getMetadataURL = func(name string) string { return server.URL }
	getDunnerTaskURLOfRecipe = func(string) string { return server.URL }
	defer server.Close()

	err := InitWithRecipe(".test_init_dunner.yaml", "foo")

	expected := "Failed to execute pre-install command from recipe: exec: \"invalid_cmd\": executable file not found in $PATH"
	if err == nil {
		t.Errorf("Expected error %s, got %s", expected, err)
	}
	if _, err := os.Stat(taskFilePath); err == nil {
		t.Errorf("Dunner task file should be removed, but it still exists")
	}
}
