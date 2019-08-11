package dunner

import (
	"fmt"
	os_user "os/user"
	"testing"

	"github.com/leopardslab/dunner/pkg/config"
	"github.com/leopardslab/dunner/pkg/docker"
)

func TestGetDunnerUserFromStep(t *testing.T) {
	expected := "test_user"
	task := config.Task{User: expected}

	user := getDunnerUser(task)

	if user != expected {
		t.Errorf("got: %s, want: %s", user, expected)
	}
}

func TestGetDunnerUserFromUserEnv(t *testing.T) {
	user, _ := os_user.Current()
	want := user.Uid

	got := getDunnerUser(config.Task{})

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
