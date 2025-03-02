package application

import (
	"context"

	"github.com/iteplenky/employee-attendance/domain"
)

type FetcherService struct {
	repo domain.FetcherRepository
}

func NewFetcherService(repo domain.FetcherRepository) *FetcherService {
	return &FetcherService{repo: repo}
}

func (s *FetcherService) GetAllAttendanceRecords(ctx context.Context) ([]domain.AttendanceEvent, error) {
	return s.repo.GetAllAttendanceRecords(ctx)
}

func (s *FetcherService) GetUserAttendanceRecords(ctx context.Context, iin string) ([]domain.AttendanceEvent, error) {
	return s.repo.GetUserAttendanceRecords(ctx, iin)
}
