package logger

import(
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func init() {
	Log.Formatter = new(logrus.TextFormatter)					// default
	Log.Formatter.(*logrus.TextFormatter).FullTimestamp = true	// Enable timestamp
	Log.Formatter.(*logrus.TextFormatter).TimestampFormat = "2006-01-02 15:04:05"	// Customize timestamp format
	Log.Level = logrus.TraceLevel
	Log.Out = os.Stdout
}