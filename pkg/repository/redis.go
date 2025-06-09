package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(addr, password string, db int) (*RedisRepo, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisRepo{client: rdb}, nil
}

func (r *RedisRepo) Set(key string, value interface{}, ttlSeconds int) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisRepo) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisRepo) Del(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}
