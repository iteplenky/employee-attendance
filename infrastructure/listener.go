package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
)

var (
	ErrListenerConnectionFailed = errors.New("listener connection failed")
	ErrListenAttendanceEvents   = errors.New("unable to listen for attendance events")
)

const listenerChannel = "attendance_events"

type Listener struct {
	listener *pq.Listener
	db       *sql.DB
	cache    *RedisCache
}

func NewListener(dbURL string, cache *RedisCache) (*Listener, error) {
	listener := pq.NewListener(dbURL, 10*time.Second, time.Minute, nil)
	if listener == nil {
		return nil, ErrListenerConnectionFailed
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = listener.Listen(listenerChannel); err != nil {
		return nil, ErrListenAttendanceEvents
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
			if err := l.cache.Publish(ctx, listenerChannel, notification.Extra); err != nil {
				log.Printf("Failed to publish event: %v\n", err)
			}
		case <-ctx.Done():
			log.Println("Shutting down listener")
			err := l.listener.Close()
			if err != nil {
				log.Println("Failed to close listener")
			}
			return
		}
	}
}

func (l *Listener) Close() error {
	if err := l.listener.Close(); err != nil {
		return err
	}
	return l.db.Close()
}
