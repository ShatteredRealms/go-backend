package config

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Version = "v1.0.0"
)

type GlobalConfig struct {
	Character   CharacterServer   `yaml:"character"`
	GameBackend GamebackendServer `yaml:"gamebackend"`
	Chat        ChatServer        `json:"chat" yaml:"chat"`
	// Uptrace     UptraceConfig     `json:"uptrace" yaml:"uptrace"`
	OpenTelemetry OpenTelemetryConfig `json:"otel" yaml:"otel"`
	Agones        AgonesConfig        `json:"agones"`
	Keycloak      KeycloakGlobal      `yaml:"keycloak"`
	Version       string
}

type KeycloakGlobal struct {
	BaseURL string `yaml:"baseURL"`
	Realm   string `yaml:"realm"`
}

func NewGlobalConfig(ctx context.Context) *GlobalConfig {
	config := &GlobalConfig{
		Character: CharacterServer{
			SROServer: SROServer{
				Local: ServerAddress{
					Port: 8081,
					Host: "",
				},
				Remote: ServerAddress{
					Port: 8081,
					Host: "",
				},
				Mode:     LocalMode,
				LogLevel: logrus.InfoLevel,
				Keycloak: KeycloakClientConfig{
					ClientId:     model.CharactersClientId,
					ClientSecret: "**********",
					Id:           "738a426a-da91-4b16-b5fc-92d63a22eb76",
				},
			},
			Postgres: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "characters",
					Username: "postgres",
					Password: "password",
				},
			},
			Mongo: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "27017",
					Username: "mongo",
					Password: "password",
					Name:     "sro",
				},
			},
		},
		GameBackend: GamebackendServer{
			SROServer: SROServer{
				Local: ServerAddress{
					Port: 8082,
					Host: "",
				},
				Remote: ServerAddress{
					Port: 8082,
					Host: "",
				},
				Mode:     LocalMode,
				LogLevel: logrus.InfoLevel,
				Keycloak: KeycloakClientConfig{
					ClientId:     model.GamebackendClientId,
					ClientSecret: "**********",
					Id:           "c3cacba8-cd16-4a4f-bc86-367274cb7cb5",
				},
			},
			Postgres: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "gamebackend",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []DBConfig{},
			},
		},
		Chat: ChatServer{
			SROServer: SROServer{
				Local: ServerAddress{
					Port: 8180,
					Host: "",
				},
				Remote: ServerAddress{
					Port: 8180,
					Host: "",
				},
				Mode:     LocalMode,
				LogLevel: logrus.InfoLevel,
				Keycloak: KeycloakClientConfig{
					ClientId:     model.ChatClientId,
					ClientSecret: "**********",
					Id:           "4c79d4a0-a3fd-495f-b56e-eea508bb0862",
				},
			},
			Kafka: ServerAddress{
				Port: 29092,
				Host: "localhost",
			},
			Postgres: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "chat",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []DBConfig{},
			},
		},
		OpenTelemetry: OpenTelemetryConfig{
			Addr: "otel-collector:4317",
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
		Keycloak: KeycloakGlobal{
			BaseURL: "http://localhost:80801/",
			Realm:   "default",
		},
		Version: Version,
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/sro/")
	viper.AddConfigPath("./test/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			log.Logger.WithContext(ctx).Fatalf("read app config parse error: %v", err)
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Logger.WithContext(ctx).Warnf("no config found, using default: %v", err)
		} else {
			log.Logger.WithContext(ctx).Fatalf("unknown error prasing config : %v", err)
		}
	}

	viper.SetEnvPrefix("SRO")
	// Read from environment variables
	helpers.BindEnvsToStruct(config)

	// Save to struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Logger.WithContext(ctx).Fatalf("unmarshal appConfig: %v", err)
	}

	return config
}
