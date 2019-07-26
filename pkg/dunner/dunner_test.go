package dunner

import (
	os_user "os/user"
	"testing"

	"github.com/leopardslab/dunner/pkg/config"
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
