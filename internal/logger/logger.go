package logger

import(
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func init() {
	Log.Formatter = new(logrus.TextFormatter)                     // default
	Log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	Log.Level = logrus.TraceLevel
	Log.Out = os.Stdout
}