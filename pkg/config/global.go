package config

import (
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Version = "v1.0.0"
)

type GlobalConfig struct {
	Characters  SROServer     `yaml:"characters"`
	GameBackend SROServer     `yaml:"gamebackend"`
	Chat        SROServer     `yaml:"chat"`
	Uptrace     UptraceConfig `yaml:"uptrace"`
	Agones      AgonesConfig  `json:"agones"`
	Version     string
}

func NewGlobalConfig() *GlobalConfig {
	config := &GlobalConfig{
		Characters: SROServer{
			Local: ServerAddress{
				Port: 8081,
				Host: "",
			},
			Remote: ServerAddress{
				Port: 8081,
				Host: "",
			},
			Mode:     LocalMode,
			LogLevel: log.InfoLevel,
			DB: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "characters",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []DBConfig{},
			},
			Kafka: ServerAddress{
				Port: 29092,
				Host: "localhost",
			},
			Keycloak: KeycloakClientConfig{
				Realm:        "default",
				BaseURL:      "http://localhost:8080",
				ClientId:     model.CharactersClientId,
				ClientSecret: "**********",
				Id:           "738a426a-da91-4b16-b5fc-92d63a22eb76",
			},
		},
		GameBackend: SROServer{
			Local: ServerAddress{
				Port: 8082,
				Host: "",
			},
			Remote: ServerAddress{
				Port: 8082,
				Host: "",
			},
			Mode:     LocalMode,
			LogLevel: log.InfoLevel,
			DB: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "gamebackend",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []DBConfig{},
			},
			Kafka: ServerAddress{
				Port: 29092,
				Host: "localhost",
			},
			Keycloak: KeycloakClientConfig{
				Realm:        "default",
				BaseURL:      "http://localhost:8080",
				ClientId:     model.GamebackendClientId,
				ClientSecret: "**********",
				Id:           "c3cacba8-cd16-4a4f-bc86-367274cb7cb5",
			},
		},
		Chat: SROServer{
			Local: ServerAddress{
				Port: 8180,
				Host: "",
			},
			Remote: ServerAddress{
				Port: 8180,
				Host: "",
			},
			Kafka: ServerAddress{
				Port: 29092,
				Host: "localhost",
			},
			Mode:     LocalMode,
			LogLevel: log.InfoLevel,
			DB: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "chat",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []DBConfig{},
			},
			Keycloak: KeycloakClientConfig{
				Realm:        "default",
				BaseURL:      "http://localhost:8080",
				ClientId:     model.ChatClientId,
				ClientSecret: "**********",
				Id:           "4c79d4a0-a3fd-495f-b56e-eea508bb0862",
			},
		},
		Uptrace: UptraceConfig{
			Host:  "localhost",
			Port:  14317,
			Id:    "2",
			Token: "project2_secret_token",
		},
		Agones: AgonesConfig{
			KeyFile:    "/etc/sro/auth/agones/client/key",
			CertFile:   "/etc/sro/auth/agones/client/cert",
			CaCertFile: "/etc/sro/auth/agones/ca/ca",
			Namespace:  "default",
			Allocator: ServerAddress{
				Port: 443,
				Host: "localhost",
			},
		},
		Version: Version,
	}

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
	helpers.BindEnvsToStruct(config)

	// Save to struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unmarshal appConfig: %v", err)
	}

	return config
}
