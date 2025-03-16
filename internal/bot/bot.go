package bot

import (
	"context"
	"errors"
	"github.com/iteplenky/employee-attendance/internal/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/iteplenky/employee-attendance/application"
)

var (
	ErrLoadSubsCache   = errors.New("failed to load subscribers to cache")
	ErrStartPolling    = errors.New("failed to start polling")
	ErrEnvTokenIsEmpty = errors.New("TOKEN environment variable not set")
	ErrTokenVerify     = errors.New("failed to verify token")
)

type Bot struct {
	Bot                 *gotgbot.Bot
	Dispatcher          *ext.Dispatcher
	UserService         *application.UserService
	SubscriptionService *application.SubscriptionService
	FetcherService      *application.FetcherService
}

func NewBot(userService *application.UserService, subs *application.SubscriptionService, fetcher *application.FetcherService) (*Bot, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, ErrEnvTokenIsEmpty
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		return nil, ErrTokenVerify
	}

	dispatcher := ext.NewDispatcher(nil)
	registerHandlers(dispatcher, userService, subs, fetcher)

	return &Bot{
		Bot:                 b,
		Dispatcher:          dispatcher,
		UserService:         userService,
		SubscriptionService: subs,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {

	err := b.SubscriptionService.LoadSubscribersToCache(ctx)
	if err != nil {
		log.Printf("failed to load subscribers to cache: %v", err)
		return ErrLoadSubsCache
	}

	updater := ext.NewUpdater(b.Dispatcher, nil)
	if err = StartPolling(b, updater); err != nil {
		log.Printf("failed to start polling: %v", err)
		return ErrStartPolling
	}

	log.Printf("%s has been started...\n", b.Bot.User.Username)
	go updater.Idle()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")

	if err = updater.Stop(); err != nil {
		log.Printf("failed to stop updater: %v", err)
	}
	log.Println("Bot stopped gracefully")

	return nil
}

func registerHandlers(dispatcher *ext.Dispatcher, userService *application.UserService, subs *application.SubscriptionService, fetcher *application.FetcherService) {
	dispatcher.AddHandler(handlers.StartHandler(userService))
	dispatcher.AddHandler(handlers.IINHandler(userService))
	dispatcher.AddHandler(handlers.ProfileCallbackHandler(userService))
	dispatcher.AddHandler(handlers.AttendanceCallbackHandler(userService, fetcher))
	dispatcher.AddHandler(handlers.NotificationsSettingsHandler(userService))
	dispatcher.AddHandler(handlers.ToggleNotificationsHandler(userService, subs))
	dispatcher.AddHandler(handlers.SettingsMenuCallbackHandler(userService))
}

func StartPolling(bot *Bot, updater *ext.Updater) error {
	return updater.StartPolling(bot.Bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout:     9,
			RequestOpts: &gotgbot.RequestOpts{Timeout: time.Second * 10},
		},
	})
}
