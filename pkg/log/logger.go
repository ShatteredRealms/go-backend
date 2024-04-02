package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.AddHook(otellogrus.NewHook(
		otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel)))

	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.TraceLevel)
	Logger.Formatter = new(prefixed.TextFormatter)
}
