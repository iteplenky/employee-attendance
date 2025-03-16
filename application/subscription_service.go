package application

import (
	"context"
	"fmt"
	"github.com/iteplenky/employee-attendance/domain"
	"strconv"
)

type SubscriptionService struct {
	db    domain.UserRepository
	cache domain.Cache
}

func NewSubscriptionService(db domain.UserRepository, cache domain.Cache) *SubscriptionService {
	return &SubscriptionService{db: db, cache: cache}
}

func (s *SubscriptionService) GetAllSubscribers(ctx context.Context) (map[string]int64, error) {
	return s.db.GetAllSubscribers(ctx)
}

func (s *SubscriptionService) LoadSubscribersToCache(ctx context.Context) error {
	subscribers, err := s.GetAllSubscribers(ctx)
	if err != nil {
		return err
	}

	for iin, userID := range subscribers {
		err = s.cache.Set(ctx, iin, strconv.FormatInt(userID, 10))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SubscriptionService) SaveSubscribersToCache(ctx context.Context, subscribers map[string]int64) error {
	for iin, userID := range subscribers {
		if err := s.cache.Set(ctx, iin, strconv.FormatInt(userID, 10)); err != nil {
			return fmt.Errorf("error adding subscriber %v: %w", userID, err)
		}
	}
	return nil
}

func (s *SubscriptionService) RemoveSubscriberFromCache(ctx context.Context, iin string) error {
	return s.cache.Del(ctx, iin)
}
