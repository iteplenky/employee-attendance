package infrastructure

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("connected to redis successfully")
	return &RedisCache{client: rdb}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Publish(ctx context.Context, channel, message string) error {
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *RedisCache) Subscribe(ctx context.Context, channel string) <-chan string {
	sub := r.client.Subscribe(ctx, channel)
	ch := make(chan string)

	go func() {
		for msg := range sub.Channel() {
			ch <- msg.Payload
		}
		close(ch)
	}()

	return ch
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
