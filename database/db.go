package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB(dbURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDB{db}, nil
}

func (p *PostgresDB) UserExists(userID int64) (bool, error) {
	var exists bool
	err := p.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE tg_id=$1)", userID).Scan(&exists)
	return exists, err
}

func (p *PostgresDB) RegisterUser(userID int64, iin string) error {
	_, err := p.Exec("INSERT INTO users (tg_id, iin) VALUES ($1, $2)", userID, iin)
	return err
}

func (p *PostgresDB) GetUser(userID int64) (*User, error) {
	var user User
	err := p.QueryRow("SELECT tg_id, iin, notifications_enabled FROM users WHERE tg_id=$1",
		userID).Scan(&user.TgID, &user.IIN, &user.NotificationsEnabled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (p *PostgresDB) EnableNotifications(userID int64) error {
	_, err := p.Exec("UPDATE users SET notifications_enabled = TRUE WHERE tg_id = $1", userID)
	return err
}

func (p *PostgresDB) AreNotificationsEnabled(userID int64) (bool, error) {
	var enabled bool
	err := p.QueryRow("SELECT notifications_enabled FROM users WHERE tg_id = $1", userID).Scan(&enabled)
	return enabled, err
}

func (p *PostgresDB) ToggleNotifications(userID int64, enabled bool) error {
	_, err := p.Exec("UPDATE users SET notifications_enabled = $1 WHERE tg_id = $2", enabled, userID)
	return err
}

func (p *PostgresDB) GetAllSubscribers() (map[string]int64, error) {
	subscribers := make(map[string]int64)

	rows, err := p.Query("SELECT iin, tg_id FROM users WHERE notifications_enabled = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var empCode string
		var tgID int64
		if err = rows.Scan(&empCode, &tgID); err != nil {
			return nil, err
		}
		subscribers[empCode] = tgID
	}

	return subscribers, nil
}
