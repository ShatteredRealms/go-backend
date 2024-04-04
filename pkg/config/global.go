package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ModeProduction  ServerMode = "production"
	ModeDevelopment ServerMode = "development"
	ModeDebug       ServerMode = "debug"
	LocalMode       ServerMode = "local"
)

type ServerMode string

var (
	Version = "v1.0.0"
)

type GlobalConfig struct {
	Character     CharacterServer     `yaml:"character"`
	GameBackend   GamebackendServer   `yaml:"gamebackend"`
	Chat          ChatServer          `json:"chat" yaml:"chat"`
	OpenTelemetry OpenTelemetryConfig `json:"otel" yaml:"otel"`
	Agones        AgonesConfig        `json:"agones"`
	Keycloak      KeycloakGlobal      `yaml:"keycloak"`
	Version       string
}

type SROServer struct {
	Local    ServerAddress        `yaml:"local"`
	Remote   ServerAddress        `yaml:"remote"`
	Mode     ServerMode           `yaml:"mode"`
	LogLevel logrus.Level         `yaml:"logLevel"`
	Keycloak KeycloakClientConfig `yaml:"keycloak"`
}

type CharacterServer struct {
	SROServer `yaml:",inline" mapstructure:",squash"`
	Postgres  DBPoolConfig `yaml:"postgres"`
	Mongo     DBPoolConfig `yaml:"mongo"`
}

type GamebackendServer struct {
	SROServer `yaml:",inline" mapstructure:",squash"`
	Postgres  DBPoolConfig `yaml:"postgres"`
}

type ChatServer struct {
	SROServer `yaml:",inline" mapstructure:",squash"`
	Postgres  DBPoolConfig  `yaml:"postgres"`
	Kafka     ServerAddress `yaml:"kafka"`
}

type KeycloakGlobal struct {
	BaseURL string `yaml:"baseURL"`
	Realm   string `yaml:"realm"`
}

// KeycloakClientConfig oidc client for keycloak
type KeycloakClientConfig struct {
	Id           string `yaml:"id"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type AgonesConfig struct {
	KeyFile    string        `yaml:"keyFile"`
	CertFile   string        `yaml:"certFile"`
	CaCertFile string        `yaml:"caCertFile"`
	Namespace  string        `yaml:"namespace"`
	Allocator  ServerAddress `yaml:"allocator"`
}

type ServerAddress struct {
	Port uint   `yaml:"port"`
	Host string `yaml:"host"`
}

func NewGlobalConfig(ctx context.Context) (*GlobalConfig, error) {
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
					ClientId:     "sro-character",
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
					ClientId:     "sro-gamebackend",
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
					ClientId:     "sro-chat",
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
			log.Logger.WithContext(ctx).Errorf("read app config parse error: %v", err)
			return nil, err
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Logger.WithContext(ctx).Infof("Using default config: %v", err)
		} else {
			log.Logger.WithContext(ctx).Errorf("unknown error prasing config : %v", err)
			return nil, err
		}
	}

	viper.SetEnvPrefix("SRO")
	// Read from environment variables
	BindEnvsToStruct(config)

	// Save to struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Logger.WithContext(ctx).Errorf("unmarshal appConfig: %v", err)
		return nil, err
	}

	return config, nil
}

func BindEnvsToStruct(obj interface{}) {
	viper.AutomaticEnv()

	val := reflect.ValueOf(obj)
	if reflect.ValueOf(obj).Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		key := field.Name
		if field.Anonymous {
			key = ""
		}
		bindRecursive(key, val.Field(i))
	}
}

func bindRecursive(key string, val reflect.Value) {
	if val.Kind() != reflect.Struct {
		env := "SRO_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		viper.MustBindEnv(key, env)
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		newKey := field.Name
		if field.Anonymous {
			newKey = ""
		} else if key != "" {
			newKey = "." + newKey
		}

		bindRecursive(key+newKey, val.Field(i))
	}
}

func (s *ServerAddress) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
