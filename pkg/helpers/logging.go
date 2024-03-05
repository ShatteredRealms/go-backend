package helpers

import (
	"fmt"
	"os"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

func SetupLogger(serviceName string) {
	log.Logger.AddHook(otellogrus.NewHook(
		otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel)))

	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(logrus.TraceLevel)
	log.Logger.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       fmt.Sprintf("%%time%% %s [%%lvl%%]: %%msg%%\n", serviceName),
	})
}
