package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
