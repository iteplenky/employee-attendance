package database

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Publish(ctx context.Context, channel, message string) error
	Subscribe(ctx context.Context, channel string) <-chan string
	Close() error
	HSet(ctx context.Context, key, field, value string) error
	HDel(ctx context.Context, key string, fields ...string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HGet(ctx context.Context, key, field string) (string, error)
}
