package main

import (
	"context"
	"github.com/iteplenky/employee-attendance/internal/bot"
	"log"

	"github.com/iteplenky/employee-attendance/application"
	"github.com/iteplenky/employee-attendance/config"
	"github.com/iteplenky/employee-attendance/infrastructure"
)

func main() {
	cfg := config.Load()

	db, err := infrastructure.NewPostgresDB(cfg.DBConnURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	cache, err := infrastructure.NewRedisCache(cfg.RedisAddr, "", 0)
	if err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	defer cache.Close()

	userService := application.NewUserService(db)
	subscriptionService := application.NewSubscriptionService(db, cache)

	b, err := bot.NewBot(userService, subscriptionService)
	if err != nil {
		log.Fatal("failed to initialize bot:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go infrastructure.HandleAttendanceEvents(ctx, cache, b)

	listener, err := infrastructure.NewListener(cfg.DBAttendanceURL, cache)
	if err != nil {
		log.Fatal("failed to initialize listener:", err)
	}
	defer listener.Close()

	go listener.StartListening(ctx)

	if err = b.Start(ctx); err != nil {
		log.Fatal("failed to start bot:", err)
	}
}
