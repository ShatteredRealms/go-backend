package helpers

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"os"
)

func SetupLogger() {
	log.AddHook(otellogrus.NewHook(
		otellogrus.WithLevels(
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel)))

	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% [%lvl%]: %msg%\n",
	})
}
