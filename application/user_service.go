package application

import (
	"context"
	"errors"

	"github.com/iteplenky/employee-attendance/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, userID int64, iin string) error {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if user != nil {
		return ErrUserAlreadyExists
	}

	return s.repo.RegisterUser(ctx, userID, iin)
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *UserService) AreNotificationsEnabled(ctx context.Context, userID int64) (bool, error) {
	return s.repo.NotificationsEnabled(ctx, userID)
}

func (s *UserService) ToggleNotifications(ctx context.Context, userID int64, enabled bool) error {
	return s.repo.ToggleNotifications(ctx, userID, enabled)
}
