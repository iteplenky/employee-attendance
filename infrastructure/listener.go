package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
)

type Listener struct {
	listener *pq.Listener
	db       *sql.DB
	cache    *RedisCache
}

func NewListener(dbURL string, cache *RedisCache) (*Listener, error) {
	listener := pq.NewListener(dbURL, 10*time.Second, time.Minute, nil)
	if listener == nil {
		return nil, errors.New("unable to connect to database listener")
	}

	if err := listener.Listen("attendance_events"); err != nil {
		return nil, errors.New("unable to listen attendance_events")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Listener{listener: listener, db: db, cache: cache}, nil
}

func (l *Listener) StartListening(ctx context.Context) {
	for {
		select {
		case notification := <-l.listener.Notify:
			if notification == nil {
				continue
			}
			log.Printf("New attendance event: %s\n", notification.Extra)
			if err := l.cache.Publish(ctx, "attendance_events", notification.Extra); err != nil {
				log.Printf("Failed to publish event: %v\n", err)
			}
		case <-ctx.Done():
			log.Println("Shutting down listener")
			l.listener.Close()
			return
		}
	}
}

func (l *Listener) Close() error {
	l.listener.Close()
	return l.db.Close()
}
