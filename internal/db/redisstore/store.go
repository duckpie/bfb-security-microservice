package redisstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/duckpie/bfb-security-microservice/internal/config"
	"github.com/go-redis/redis/v8"
)

type redisStore struct {
	client *redis.Client
}

type RedisStoreI interface {
	Save(ctx context.Context, key string, value interface{}, d time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	DbSize(ctx context.Context) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Delete(ctx context.Context, keys ...string) error
	Clear(ctx context.Context) error
	Close() error
}

// Сохранить объект в БД
func (s *redisStore) Save(ctx context.Context, key string, value interface{}, d time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, b, d).Err()
}

// Получить данные из хранилища
func (s *redisStore) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

// Выборка всех ключей по заданному паттерну
func (s *redisStore) Keys(ctx context.Context, pattern string) ([]string, error) {
	return s.client.Keys(ctx, pattern).Result()
}

// Получить кол-во записей в хранилище
func (s *redisStore) DbSize(ctx context.Context) (int64, error) {
	return s.client.DBSize(ctx).Result()
}

// Удалить по переданным ключам записи
func (s *redisStore) Delete(ctx context.Context, keys ...string) error {
	return s.client.Del(ctx, keys...).Err()
}

// Закрыть соединение
func (s *redisStore) Close() error {
	return s.client.Close()
}

// Очистить хранилище
func (s *redisStore) Clear(ctx context.Context) error {
	return s.client.FlushAllAsync(ctx).Err()
}

func NewClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:              int(cfg.DB),
		PoolSize:        5000,
		MinRetryBackoff: 1 * time.Second,
		MaxRetryBackoff: 2 * time.Second,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisStore(c *redis.Client) RedisStoreI {
	return &redisStore{
		client: c,
	}
}
