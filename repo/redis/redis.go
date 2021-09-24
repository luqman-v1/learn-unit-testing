package redis

import (
	"context"
	"time"
)

type IRedis interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, interface{}, time.Duration) error
	Del(ctx context.Context, key string) error
}

type RedisDB struct{}

func (r RedisDB) Get(ctx context.Context, s string) ([]byte, error) {
	panic("implement me")
}

func (r RedisDB) Set(ctx context.Context, s string, i interface{}, duration time.Duration) error {
	panic("implement me")
}

func (r RedisDB) Del(ctx context.Context, key string) error {
	panic("implement me")
}

func NewRedisDB() *RedisDB {
	return &RedisDB{}
}
