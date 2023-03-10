package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"os"
	"reflect"
	"strings"
)

func SetupConfig(conf interface{}) {
	// Load file appConfig
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./test/")
	viper.AddConfigPath("/etc/sro/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			log.Fatalf("read appConfig: %v", err)
		}
	}

	// Read from environment variables
	viper.SetEnvPrefix("SRO")
	BindEnvsToStruct(conf)

	// Save to struct
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("unmarshal appConfig: %v", err)
	}
}

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

func BindEnvsToStruct(obj interface{}) {
	viper.AutomaticEnv()

	val := reflect.ValueOf(obj)
	if reflect.ValueOf(obj).Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		bindRecursive(val.Type().Field(i).Name, val.Field(i))
	}
}

func bindRecursive(name string, val reflect.Value) {
	if val.Kind() != reflect.Struct {
		viper.BindEnv(name, strings.ReplaceAll(name, ".", "_"))
		return
	}

	for i := 0; i < val.NumField(); i++ {
		bindRecursive(name+"."+val.Type().Field(i).Name, val.Field(i))
	}
}
