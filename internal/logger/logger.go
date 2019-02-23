package logger

import (
	"os"

	"github.com/sirupsen/logrus"
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
