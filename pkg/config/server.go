package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	ModeProduction  ServerMode = "production"
	ModeDevelopment ServerMode = "development"
	ModeDebug       ServerMode = "debug"
	LocalMode       ServerMode = "local"
)

type ServerMode string

type SROServer struct {
	Local    ServerAddress        `yaml:"local"`
	Remote   ServerAddress        `yaml:"remote"`
	Mode     ServerMode           `yaml:"mode"`
	LogLevel log.Level            `yaml:"logLevel"`
	Keycloak KeycloakClientConfig `yaml:"keycloak"`
}

type CharacterServer struct {
	SROServer `yaml:"server"`
	Postgres  DBPoolConfig `yaml:"postgres"`
	Mongo     DBPoolConfig `yaml:"mongo"`
}

type GamebackendServer struct {
	SROServer `yaml:"server"`
	Postgres  DBPoolConfig `yaml:"postgres"`
}

type ChatServer struct {
	SROServer `yaml:"server"`
	Postgres  DBPoolConfig  `yaml:"postgres"`
	Kafka     ServerAddress `yaml:"kafka"`
}

type ServerAddress struct {
	Port uint   `yaml:"port"`
	Host string `yaml:"host"`
}

func (s *ServerAddress) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
