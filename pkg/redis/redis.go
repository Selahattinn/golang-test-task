package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	connTimeout = 20 * time.Second
)

type Config struct {
	URL      string
	Port     string
	Password string
	Database int
	TLS      bool
}

type IRedisInstance interface {
	Close() error
	Get(ctx context.Context, key string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Ping(ctx context.Context) error
}

type redisInstance struct {
	client *redis.Client
}

func InitRedis(config Config) (IRedisInstance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.URL, config.Port),
		Password: config.Password,
		DB:       config.Database,
	}
	if config.TLS {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisInstance{client: rdb}, nil
}

func (r *redisInstance) Close() error {
	return r.client.Close()
}

func (r *redisInstance) Get(ctx context.Context, key string) ([]string, error) {
	//find keys with prefix
	keys, err := r.client.Keys(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	//get values for keys
	var messages []string
	for _, key := range keys {
		fmt.Println("key", key)
		val, err := r.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		messages = append(messages, val)
	}

	return messages, nil
}

func (r *redisInstance) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisInstance) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
func (r *redisInstance) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
