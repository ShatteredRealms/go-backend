package repository

import (
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// ConnectDB Initializes the connection to a Postgres database
func ConnectDB(pool config.DBPoolConfig) (*gorm.DB, error) {
	conf, err := pgx.ParseConfig(pool.Master.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	sqlDB := stdlib.OpenDB(*conf)

	sqlDB.SetConnMaxLifetime(time.Second)
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetMaxIdleConns(10)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.New(
			log.Logger,
			logger.Config{
				SlowThreshold:             time.Millisecond * 500,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
			},
		),
	})

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

// func ConnectFromFile(filePath string) (*gorm.DB, error) {
// 	file, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	c := &config.DBPoolConfig{}
// 	err = yaml.Unmarshal(file, c)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return ConnectDB(*c)
// }
