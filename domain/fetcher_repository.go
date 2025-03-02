package domain

import "context"

type FetcherRepository interface {
	GetAllAttendanceRecords(ctx context.Context) ([]AttendanceEvent, error)
	GetUserAttendanceRecords(ctx context.Context, iin string) ([]AttendanceEvent, error)
}
