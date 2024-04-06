package config

import "fmt"

// DBConfig Information on how to connect to the database
type DBConfig struct {
	ServerAddress `yaml:",inline" mapstructure:",squash"`
	Name          string `yaml:"name"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
}

// DBPoolConfig Defines the master and slave connections to a replicated database. Slaves may be empty.
type DBPoolConfig struct {
	Master DBConfig   `yaml:"master" json:"master"`
	Slaves []DBConfig `yaml:"slaves" json:"slaves"`
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

func (c DBConfig) MongoDSN() string {
	// refer https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/ for details
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?timeoutMS=5000",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
	)
}

func (p DBPoolConfig) Addresses() []string {
	addrs := make([]string, len(p.Slaves)+1)
	addrs[0] = p.Master.Address()
	for idx, dbConf := range p.Slaves {
		addrs[idx+1] = dbConf.Address()
	}

	return addrs
}
