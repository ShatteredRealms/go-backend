package config

import (
	log "github.com/sirupsen/logrus"
)

const (
	ModeProduction  ServerMode = "production"
	ModeDevelopment ServerMode = "development"
	ModeDebug       ServerMode = "debug"
	LocalMode       ServerMode = "local"
)

var (
	AllModes = []ServerMode{ModeProduction, ModeDebug, ModeDevelopment}
)

type ServerMode string

type SROServer struct {
	Local    ServerAddress        `yaml:"local"`
	Remote   ServerAddress        `yaml:"remote"`
	Mode     ServerMode           `yaml:"mode"`
	LogLevel log.Level            `yaml:"logLevel"`
	DB       DBPoolConfig         `yaml:"db"`
	Kafka    ServerAddress        `yaml:"kafka"`
	Keycloak KeycloakClientConfig `yaml:"keycloak"`
}

type ServerAddress struct {
	Port uint   `yaml:"port"`
	Host string `yaml:"host"`
}
