package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Log is a globally configured logger
var Log = logrus.New()

func init() {
	Log.Formatter = new(logrus.TextFormatter)                                     // Default
	Log.Formatter.(*logrus.TextFormatter).FullTimestamp = true                    // Enable timestamp
	Log.Formatter.(*logrus.TextFormatter).TimestampFormat = "2006-01-02 15:04:05" // Customize timestamp format
	Log.Level = logrus.TraceLevel
	Log.Out = os.Stdout
}

// InitColorOutput disables colorized output if no-color flag is passed
func InitColorOutput() {
	if viper.GetBool("No-color") {
		color.NoColor = true
	}
}

// ErrorOutput prints the given message in red color
func ErrorOutput(format string, a ...interface{}) {
	color.Red(format, a...)
}

// Bullet prints out the given message into stdout with a bulleted symbol at start
func Bullet(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("• "+format, a...))
}
