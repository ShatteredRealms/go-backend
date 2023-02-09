package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	ModeProduction  ServerMode = "production"
	ModeDebug       ServerMode = "debug"
	ModeVerbose     ServerMode = "verbose"
	ModeDevelopment ServerMode = "development"
)

var (
	AllModes = []ServerMode{ModeProduction, ModeDebug, ModeVerbose, ModeDevelopment}
)

type ServerMode string

type Server struct {
	Local    ServerAddress `yaml:"local"`
	Remote   ServerAddress `yaml:"remote"`
	Mode     ServerMode    `yaml:"mode"`
	LogLevel log.Level     `yaml:"logLevel"`
	DB       DBPoolConfig  `yaml:"db"`
	Kafka    ServerAddress `yaml:"kafka"`
}

type ServerAddress struct {
	Port uint   `yaml:"port"`
	Host string `yaml:"host"`
}

// DBConfig Information on how to connect to the database
type DBConfig struct {
	Host     string `yaml:"hoster"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// DBPoolConfig Defines the master and slave connections to a replicated database. Slaves may be empty.
type DBPoolConfig struct {
	Master DBConfig   `yaml:"master"`
	Slaves []DBConfig `yaml:"slaves"`
}

func (s *ServerAddress) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (c DBConfig) MySQLDSN() string {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

func (c DBConfig) PostgresDSN() string {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}
