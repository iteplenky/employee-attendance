package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type Listener struct {
	listener *pq.Listener
	db       *sql.DB
}

func NewListener(dbURL string) (*Listener, error) {

	listener := pq.NewListener(dbURL, 10*time.Second, time.Minute, nil)
	if listener == nil {
		log.Println("unable to connect to database", listener)
		return nil, errors.New("unable to connect to database")
	}

	if err := listener.Listen("attendance_events"); err != nil {
		fmt.Printf("could not listen attendance_events: %v\n", err)
		return nil, errors.New("unable to listen attendance_events")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("unable to connect to postgres: %v\n", err)
		return nil, errors.New("unable to connect to postgres")
	}

	if err = db.Ping(); err != nil {
		fmt.Printf("unable to ping postgres: %v\n", err)
		return nil, errors.New("unable to ping postgres")
	}

	log.Printf("connected to attendance_events\n")
	return &Listener{listener: listener, db: db}, nil
}

func (l *Listener) StartListening() {
	log.Printf("listening attendance_events\n")

	for {
		select {
		case notification := <-l.listener.Notify:
			if notification == nil {
				continue
			}
			log.Printf("notification: %v\n", notification.Extra)
		}
	}
}

func (l *Listener) Close() {
	err := l.listener.Close()
	if err != nil {
		log.Printf("listener close error: %v\n", err)
		return
	}
	err = l.db.Close()
	if err != nil {
		log.Printf("db close error: %v\n", err)
		return
	}
}
