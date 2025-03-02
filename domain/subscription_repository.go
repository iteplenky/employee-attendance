package domain

import "context"

type SubscriptionRepository interface {
	Subscribe(ctx context.Context, sub Subscription) error
	Unsubscribe(ctx context.Context, userID int64) error
}
