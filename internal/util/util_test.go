package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

func TestDirExistsSuccess(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "TestDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	exists := DirExists(tmpdir)

	if !exists {
		t.Fatalf("Directory exists; but got false")
	}
}

func TestDirExistsFail(t *testing.T) {
	exists := DirExists("this path is invalid")

	if exists {
		t.Fatalf("Directory invalid; but got as exists")
	}
}

func TestDirExistsFailForFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestFileExists")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpfile.Name())

	exists := DirExists(tmpfile.Name())

	if exists {
		t.Fatalf("Not a directory; but got as true")
	}
}

func TestDirExistsIfNotAbsPath(t *testing.T) {
	exists := DirExists("~/invalidpathfortesting")

	if exists {
		t.Fatalf("Not a directory; but got as true")
	}
}

func TestDownloadFailureIfFileNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	err := Download(server.URL, ".test_dunner.yaml")

	if err == nil {
		t.Fatalf("Expected %s, got nil", "error")
	}
}

func TestDownloadFailureIfServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	err := Download(server.URL, ".test_dunner.yaml")

	if err == nil {
		t.Fatalf("Expected %s, got nil", "error")
	}
}

func TestDownloadFailureIfSomeHTTPError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	err := Download(server.URL, ".test_dunner.yaml")

	if err == nil {
		t.Fatalf("Expected %s, got nil", "error")
	}
}

func TestDownloadSuccess(t *testing.T) {
	defer os.Remove(".test_dunner.yaml")

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "All OK", http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	err := Download(server.URL, ".test_dunner.yaml")
	expectedFileContents := "All OK\n"

	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}

	fileContents, err := ioutil.ReadFile(".test_dunner.yaml")
	if err != nil {
		t.Fatalf("Failed to read downloaded dunner yaml file")
	}
	if string(fileContents) != expectedFileContents {
		t.Fatalf("Downloaded file contents not matching. Expected %s, got %s", expectedFileContents, string(fileContents))
	}
}

func TestGetUrlContents(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "All OK", http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	fileContents, err := GetURLContents(server.URL)
	expectedFileContents := "All OK\n"

	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}
	if string(fileContents) != expectedFileContents {
		t.Fatalf("Downloaded file contents not matching. Expected %s, got %s", expectedFileContents, string(fileContents))
	}
}

func TestGetUrlContents404(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	fileContents, err := GetURLContents(server.URL)

	expected := fmt.Sprintf("Error downloading file %s: 404 Not Found", server.URL)
	if fileContents != nil {
		t.Fatalf("Expected no fileContents, got %s", fileContents)
	}
	if err == nil || err.Error() != expected {
		t.Fatalf("Expected error %s, got %s", expected, err)
	}
}

func TestGetExecutableCommandForSystemCommands(t *testing.T) {
	wd, _ := os.Getwd()

	command := "go"
	cmd := getExecutableCommand(command)

	pd, _ := os.Getwd()
	expectedCommand := exec.Command("go")

	if wd != pd {
		t.Errorf("Working directory changed, expected %s, got %s", wd, pd)
	}
	if cmd == nil {
		t.Errorf("Command nil, expected valid command")
	}
	if cmd.Path != expectedCommand.Path {
		t.Errorf("Command path invalid, expected %s, got %s", expectedCommand.Path, cmd.Path)
	}
}

func TestGetExecutableCommandForCommandsWithPath(t *testing.T) {
	wd, _ := os.Getwd()
	logger1 := createLogger("logger1")
	logger2 := createLogger("logger2")
	command := "/bin/java"

	cmd := prepareCommand([]string{command, "-v", "-d"}, logger1, logger2)

	pd, _ := os.Getwd()
	args := make(map[string]bool)
	for _, v := range cmd.Args {
		args[v] = true
	}

	if wd != pd {
		t.Errorf("Working directory changed, expected %s, got %s", wd, pd)
	}
	if cmd == nil {
		t.Errorf("Command nil, expected valid command")
	}
	if cmd.Path != command {
		t.Errorf("Command path invalid, expected %s, got %s", command, cmd.Path)
	}
	if cmd.Dir != wd {
		t.Errorf("Command directory invalid")
	}
	if !logger1.equals(cmd.Stdout.(testLogger)) {
		t.Errorf("Logger1 invalid")
	}
	if !logger2.equals(cmd.Stderr.(testLogger)) {
		t.Errorf("Logger2 invalid")
	}
	if args["-v"] != true || args["-d"] != true {
		t.Errorf("Command arguments not valid")
	}
}

type testLogger struct {
	name string
}

func (l testLogger) Write(b []byte) (n int, err error) {
	return 1, nil
}

func createLogger(name string) testLogger {
	return testLogger{name}
}

func (l testLogger) equals(l1 testLogger) bool {
	return l.name == l1.name
}

func TestFileExists(t *testing.T) {
	if FileExists("util.go") != true {
		t.Errorf("expected file to exist, but does not exist")
	}
	if FileExists("invalid") != false {
		t.Errorf("expected file to not exist, but exists")
	}
}
