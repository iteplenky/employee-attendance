package main

import (
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/iteplenky/employee-attendance/config"
	"github.com/iteplenky/employee-attendance/database"
	"github.com/iteplenky/employee-attendance/internal/attendance"
	"github.com/iteplenky/employee-attendance/internal/bot"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	db, err := database.NewPostgresDB(cfg.DBConnURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	cache, err := database.NewRedisCache(cfg.RedisAddr, "", 0)
	if err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	defer cache.Close()

	listener, err := database.NewListener(cfg.DBAttendanceURL, cache)
	if err != nil {
		log.Fatal("failed to initialize listener:", err)
	}
	defer listener.Close()

	b := bot.NewBot(db, cache)
	subscribers, err := db.GetAllSubscribers()
	if err != nil {
		log.Fatal("failed to get all subscribers:", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listener.StartListening(ctx)

	bot.LoadSubscribersToCache(ctx, cache, subscribers)

	go attendance.HandleAttendanceEvents(ctx, cache, b)

	updater := ext.NewUpdater(b.Dispatcher, nil)
	if err = bot.StartPolling(b, updater); err != nil {
		log.Fatal("failed to start polling:", err)
	}
	log.Printf("%s has been started...\n", b.Bot.User.Username)

	go updater.Idle()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
	cancel()
	updater.Stop()
	log.Println("Bot stopped gracefully")
}
