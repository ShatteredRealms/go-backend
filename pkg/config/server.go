package config

import (
	"fmt"
)

type ServerAddress struct {
	Port uint   `yaml:"port"`
	Host string `yaml:"host"`
}

func (s *ServerAddress) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
