package config

import "fmt"

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
	// refer https://github.com/go-sql-driver/postsgres#dsn-data-source-name for details
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}
