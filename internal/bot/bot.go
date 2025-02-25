package bot

import (
	"github.com/iteplenky/employee-attendance/database"
	"github.com/iteplenky/employee-attendance/internal/handlers"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Bot struct {
	Bot        *gotgbot.Bot
	Dispatcher *ext.Dispatcher
}

func NewBot(db database.UserRepository) *Bot {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(nil)
	RegisterHandlers(dispatcher, db)

	return &Bot{Bot: b, Dispatcher: dispatcher}
}

func RegisterHandlers(dispatcher *ext.Dispatcher, db database.UserRepository) {
	dispatcher.AddHandler(handlers.StartHandler(db))
	dispatcher.AddHandler(handlers.IINHandler(db))
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
