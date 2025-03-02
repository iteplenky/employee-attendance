package domain

import "context"

type UserRepository interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
	RegisterUser(ctx context.Context, userID int64, iin string) error
	GetUser(ctx context.Context, userID int64) (*User, error)
	EnableNotifications(ctx context.Context, userID int64) error
	AreNotificationsEnabled(ctx context.Context, userID int64) (bool, error)
	ToggleNotifications(ctx context.Context, userID int64, enabled bool) error
	GetAllSubscribers(ctx context.Context) (map[string]int64, error)
	Close() error
}

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
