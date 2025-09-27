package common

import (
	"context"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

type RedisCfg struct {
	Addr         string
	Pwd          string
	Db           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
}

func GetRedis(cfg RedisCfg) *redis.Client {
	options := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pwd,
		DB:           cfg.Db,
		PoolSize:     cfg.PoolSize,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MinIdleConns,
	}
	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	client.AddHook(redisotel.NewTracingHook())
	return client
}
