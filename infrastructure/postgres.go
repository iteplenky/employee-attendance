package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"github.com/iteplenky/employee-attendance/domain"
	_ "github.com/lib/pq"
	"log"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(dbURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) UserExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := p.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM telegram_bot.users WHERE tg_id=$1)", userID).Scan(&exists)
	return exists, err
}

func (p *PostgresDB) RegisterUser(ctx context.Context, userID int64, iin string) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO telegram_bot.users (tg_id, iin) VALUES ($1, $2)", userID, iin)
	return err
}

func (p *PostgresDB) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	var user domain.User
	err := p.db.QueryRowContext(ctx, "SELECT tg_id, iin, notifications_enabled FROM telegram_bot.users WHERE tg_id=$1",
		userID).Scan(&user.ID, &user.IIN, &user.NotificationsEnabled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (p *PostgresDB) NotificationsEnabled(ctx context.Context, userID int64) (bool, error) {
	var enabled bool
	err := p.db.QueryRowContext(ctx, "SELECT notifications_enabled FROM telegram_bot.users WHERE tg_id = $1", userID).Scan(&enabled)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return enabled, err
}

func (p *PostgresDB) ToggleNotifications(ctx context.Context, userID int64, enabled bool) error {
	_, err := p.db.ExecContext(ctx, "UPDATE telegram_bot.users SET notifications_enabled = $1 WHERE tg_id = $2", enabled, userID)
	return err
}

func (p *PostgresDB) GetAllSubscribers(ctx context.Context) (map[string]int64, error) {
	subscribers := make(map[string]int64)

	rows, err := p.db.QueryContext(ctx, "SELECT iin, tg_id FROM telegram_bot.users WHERE notifications_enabled = TRUE")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var iin string
		var tgID int64
		if err = rows.Scan(&iin, &tgID); err != nil {
			return nil, err
		}
		subscribers[iin] = tgID
	}

	return subscribers, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
