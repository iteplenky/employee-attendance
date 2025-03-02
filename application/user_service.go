package application

import (
	"context"
	"errors"

	"github.com/iteplenky/employee-attendance/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, userID int64, iin string) error {
	exists, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	return s.repo.RegisterUser(ctx, userID, iin)
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *UserService) AreNotificationsEnabled(ctx context.Context, userID int64) (bool, error) {
	return s.repo.AreNotificationsEnabled(ctx, userID)
}

func (s *UserService) ToggleNotifications(ctx context.Context, userID int64, enabled bool) error {
	return s.repo.ToggleNotifications(ctx, userID, enabled)
}
