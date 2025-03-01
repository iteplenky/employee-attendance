package bot

import (
	"context"
	"github.com/iteplenky/employee-attendance/database"
	"github.com/iteplenky/employee-attendance/internal/handlers"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Bot struct {
	Bot        *gotgbot.Bot
	Dispatcher *ext.Dispatcher
}

func NewBot(db database.UserRepository, cache database.Cache) *Bot {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(nil)
	RegisterHandlers(dispatcher, db, cache)

	return &Bot{Bot: b, Dispatcher: dispatcher}
}

func RegisterHandlers(dispatcher *ext.Dispatcher, db database.UserRepository, cache database.Cache) {
	dispatcher.AddHandler(handlers.StartHandler(db))
	dispatcher.AddHandler(handlers.IINHandler(db))
	dispatcher.AddHandler(handlers.ProfileCallbackHandler(db))
	dispatcher.AddHandler(handlers.NotificationsSettingsHandler(db))
	dispatcher.AddHandler(handlers.ToggleNotificationsHandler(db, cache))
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

func LoadSubscribersToCache(ctx context.Context, cache database.Cache, subscribers map[string]int64) {
	for iin, userID := range subscribers {
		err := cache.HSet(ctx, "subscribed_users", iin, strconv.FormatInt(userID, 10))
		if err != nil {
			log.Printf("error adding subscriber %v: %v", userID, err)
		}
	}
	log.Printf("loaded %d subscribers into cache", len(subscribers))
}
