package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/fatih/color"
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
