package database

import (
	"context"
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
	cache    Cache
}

func NewListener(dbURL string, cache Cache) (*Listener, error) {

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
	return &Listener{
		listener: listener,
		db:       db,
		cache:    cache,
	}, nil
}

func (l *Listener) StartListening(ctx context.Context) {
	log.Printf("listening attendance_events\n")

	for {
		select {
		case notification := <-l.listener.Notify:
			if notification == nil {
				continue
			}
			log.Printf("new notification: %v\n", notification.Extra)
			if err := l.cache.Publish(ctx, "attendance_events", notification.Extra); err != nil {
				log.Printf("unable to publish attendance_events: %v\n", err)
			}
		case <-ctx.Done():
			log.Println("shutting down listener")
			if err := l.listener.Close(); err != nil {
				log.Printf("unable to close listener: %v\n", err)
			}
			return
		}
	}
}

func (l *Listener) Close() error {
	if err := l.listener.Close(); err != nil {
		return err
	}
	if err := l.db.Close(); err != nil {
		return err
	}
	return nil
}
