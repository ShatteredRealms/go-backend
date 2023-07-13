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
	Character   CharacterServer   `yaml:"characters"`
	GameBackend GamebackendServer `yaml:"gamebackend"`
	Chat        ChatServer        `yaml:"chat"`
	Uptrace     UptraceConfig     `yaml:"uptrace"`
	Agones      AgonesConfig      `json:"agones"`
	Version     string
}

func NewGlobalConfig() *GlobalConfig {
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
				LogLevel: log.InfoLevel,
				Keycloak: KeycloakClientConfig{
					Realm:        "default",
					BaseURL:      "http://localhost:8080",
					ClientId:     model.CharactersClientId,
					ClientSecret: "**********",
					Id:           "738a426a-da91-4b16-b5fc-92d63a22eb76",
				},
			},
			PostgresConfig: DBPoolConfig{
				Master: DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "characters",
					Username: "postgres",
					Password: "password",
				},
			},
			MongoConfig: DBPoolConfig{
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
				LogLevel: log.InfoLevel,
				Keycloak: KeycloakClientConfig{
					Realm:        "default",
					BaseURL:      "http://localhost:8080",
					ClientId:     model.GamebackendClientId,
					ClientSecret: "**********",
					Id:           "c3cacba8-cd16-4a4f-bc86-367274cb7cb5",
				},
			},
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
				LogLevel: log.InfoLevel,
				Keycloak: KeycloakClientConfig{
					Realm:        "default",
					BaseURL:      "http://localhost:8080",
					ClientId:     model.ChatClientId,
					ClientSecret: "**********",
					Id:           "4c79d4a0-a3fd-495f-b56e-eea508bb0862",
				},
			},
			Kafka: ServerAddress{
				Port: 29092,
				Host: "localhost",
			},
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
		},
		Uptrace: UptraceConfig{
			DSN: "http://project2_secret_token@localhost:14317/2",
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
