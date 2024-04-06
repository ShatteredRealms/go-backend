package cacher

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/go-gorm/caches/v4"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type redisCacher struct {
	rdb *redis.ClusterClient
}

func (c *redisCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	res, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := q.Unmarshal([]byte(res)); err != nil {
		return nil, err
	}

	return q, nil
}

func (c *redisCacher) Store(ctx context.Context, key string, val *caches.Query[any]) error {
	res, err := val.Marshal()
	if err != nil {
		return err
	}

	c.rdb.Set(ctx, key, res, 300*time.Second) // Set proper cache time
	return nil
}

func (c *redisCacher) Invalidate(ctx context.Context) error {
	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = c.rdb.Scan(ctx, cursor, fmt.Sprintf("%s*", caches.IdentifierPrefix), 0).Result()
		if err != nil {
			return err
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if _, err := c.rdb.Del(ctx, keys...).Result(); err != nil {
			log.Logger.WithContext(ctx).Infof("failed to delete %d keys, retrying indiviually: %v", len(keys), err)
			for _, key := range keys {
				if _, err := c.rdb.Del(ctx, key).Result(); err != nil {
					log.Logger.WithContext(ctx).Errorf("failed to delete key %s: %v", key, err)
					return err
				}
			}
		}
	}
	return nil
}

func NewRedisCache(ctx context.Context, dbPoolConf config.DBPoolConfig) (caches.Cacher, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: dbPoolConf.Addresses(),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			withPassword := "no password"
			if dbPoolConf.Master.Password != "" {
				withPassword = "password"
			}
			withUsername := "no username"
			if dbPoolConf.Master.Username != "" {
				withUsername = "a username"
			}
			log.Logger.WithContext(ctx).Debugf("Connecting to redis with %s and %s", withUsername, withPassword)
			_, err := cn.Ping(context.Background()).Result()
			return err
		},
		Username: dbPoolConf.Master.Username,
		Password: dbPoolConf.Master.Password,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("tracing instrumentation: %w", err)
	}

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		return nil, fmt.Errorf("metrics instrumentation: %w", err)
	}

	return &redisCacher{
		rdb: rdb,
	}, nil
}
