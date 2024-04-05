package cacher

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
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
			return err
		}
	}
	return nil
}

func NewRedisCache(dbPoolConf config.DBPoolConfig) (caches.Cacher, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    dbPoolConf.Addresses(),
		Username: dbPoolConf.Master.Username,
		Password: dbPoolConf.Master.Password,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("tracing instrumentation: %w", err)
	}

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		return nil, fmt.Errorf("metrics instrumentation: %w", err)
	}

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &redisCacher{
		rdb: rdb,
	}, nil
}
