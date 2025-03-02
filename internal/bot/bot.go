package bot

import (
	"context"
	"fmt"
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

type Bot struct {
	Bot                 *gotgbot.Bot
	Dispatcher          *ext.Dispatcher
	UserService         *application.UserService
	SubscriptionService *application.SubscriptionService
}

func NewBot(userService *application.UserService, subs *application.SubscriptionService) (*Bot, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}

	dispatcher := ext.NewDispatcher(nil)
	registerHandlers(dispatcher, userService, subs)

	return &Bot{
		Bot:                 b,
		Dispatcher:          dispatcher,
		UserService:         userService,
		SubscriptionService: subs,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	subscribers, err := b.SubscriptionService.GetAllSubscribers(ctx)
	if err != nil {
		return fmt.Errorf("get all subscribers failed: %w", err)
	}

	if err = b.SubscriptionService.SaveSubscribersToCache(ctx, subscribers); err != nil {
		return fmt.Errorf("failed to load subscribers: %w", err)
	}

	updater := ext.NewUpdater(b.Dispatcher, nil)
	if err = StartPolling(b, updater); err != nil {
		return fmt.Errorf("failed to start polling: %w", err)
	}

	log.Printf("%s has been started...\n", b.Bot.User.Username)
	go updater.Idle()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")

	updater.Stop()
	log.Println("Bot stopped gracefully")

	return nil
}

func registerHandlers(dispatcher *ext.Dispatcher, userService *application.UserService, subs *application.SubscriptionService) {
	dispatcher.AddHandler(handlers.StartHandler(userService))
	dispatcher.AddHandler(handlers.IINHandler(userService))
	dispatcher.AddHandler(handlers.ProfileCallbackHandler(userService))
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
