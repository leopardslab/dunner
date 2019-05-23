package util

import (
	"io/ioutil"
	"os"
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
