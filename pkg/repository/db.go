package repository

import (
	"fmt"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"io/ioutil"
	"time"
)

// ConnectDB Initializes the connection to a Postgres database
func ConnectDB(pool config.DBPoolConfig) (*gorm.DB, error) {
	config, err := pgx.ParseConfig(pool.Master.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	sqlDB := stdlib.OpenDB(*config)

	sqlDB.SetConnMaxLifetime(time.Second)
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetMaxIdleConns(10)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("gorm: %w", err)
	}

	if len(pool.Slaves) > 0 {
		replicas := make([]gorm.Dialector, len(pool.Slaves))
		for _, slave := range pool.Slaves {
			replicas = append(replicas, postgres.Open(slave.PostgresDSN()))
		}

		err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(pool.Master.Name))); err != nil {
		return nil, err
	}

	return db, err
}

func ConnectFromFile(filePath string) (*gorm.DB, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	c := &config.DBPoolConfig{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	return ConnectDB(*c)
}
