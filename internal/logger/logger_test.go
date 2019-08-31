package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func TestErrorOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	oldOutput := color.Output
	color.Output = buf
	name := "dunner"
	message := "Welcome!"

	ErrorOutput("Hello %s! %s", name, message)

	line, _ := buf.ReadString('\n')
	got := fmt.Sprintf("%q", line)
	expected := fmt.Sprintf("Hello %s! %s\n", name, message)
	escaped := fmt.Sprintf("%q", expected)
	color.Output = oldOutput

	if got != escaped {
		t.Fatalf("expected: %s, got: %s", escaped, got)
	}
}

func TestInitColorOutput_True(t *testing.T) {
	viper.Set("No-color", true)

	InitColorOutput()

	if color.NoColor != true {
		t.Fatalf("expected no-color to be set as true, but got %v", color.NoColor)
	}
}

func ExampleBullet() {
	arg := "foobar"

	Bullet("setup %s", arg)

	// Output: • setup foobar
}
