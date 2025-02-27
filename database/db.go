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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (p *PostgresDB) RegisterUser(userID int64, iin string) error {
	_, err := p.Exec("INSERT INTO users (tg_id, iin) VALUES ($1, $2)", userID, iin)
	return err
}

func (p *PostgresDB) GetUser(userID int64) (*User, error) {
	var user User
	err := p.QueryRow("SELECT tg_id, iin FROM users WHERE tg_id=$1", userID).Scan(&user.TgID, &user.IIN)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresDB) SaveSchedule(userID int64, startTime, endTime string) error {
	_, err := p.Exec(`
		INSERT INTO schedules (user_id, start_time, end_time) 
		VALUES ((SELECT id FROM users WHERE tg_id=$1), $2, $3) 
		ON CONFLICT (user_id) 
		DO UPDATE SET start_time=$2, end_time=$3`,
		userID, startTime, endTime)
	return err
}

func (p *PostgresDB) GetSchedule(userID int64) (string, string, error) {
	var startTime, endTime string
	err := p.QueryRow(`
		SELECT start_time, end_time FROM schedules 
		WHERE user_id = (SELECT id FROM users WHERE tg_id=$1)`,
		userID).Scan(&startTime, &endTime)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil
		}
		return "", "", err
	}
	return startTime, endTime, nil
}
