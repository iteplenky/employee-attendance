package infrastructure

import (
	"context"
	"database/sql"
	"github.com/iteplenky/employee-attendance/domain"
)

type Fetcher struct {
	db *sql.DB
}

func NewFetcher(dbURL string) (*Fetcher, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Fetcher{db: db}, nil
}

func (f *Fetcher) GetAllAttendanceRecords(ctx context.Context) ([]domain.AttendanceEvent, error) {
	rows, err := f.db.QueryContext(ctx, "SELECT id, emp_id, punch_time, terminal_alias, processed FROM attendance_log")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.AttendanceEvent
	for rows.Next() {
		var event domain.AttendanceEvent
		if err := rows.Scan(&event.ID, &event.IIN, &event.PunchTime, &event.TerminalAlias, &event.Processed); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (f *Fetcher) GetUserAttendanceRecords(ctx context.Context, iin string) ([]domain.AttendanceEvent, error) {
	rows, err := f.db.QueryContext(ctx, `
		SELECT id, emp_id, punch_time, terminal_alias, processed 
		FROM attendance_log 
		WHERE emp_id = $1
		ORDER BY punch_time DESC`, iin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.AttendanceEvent
	for rows.Next() {
		var event domain.AttendanceEvent
		if err = rows.Scan(&event.ID, &event.IIN, &event.PunchTime, &event.TerminalAlias, &event.Processed); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (f *Fetcher) Close() error {
	return f.db.Close()
}
