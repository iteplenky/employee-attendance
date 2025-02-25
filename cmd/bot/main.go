package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/iteplenky/employee-attendance/internal/bot"
	"log"
)

func main() {
	b := bot.NewBot()
	updater := ext.NewUpdater(b.Dispatcher, nil)
	if err := bot.StartPolling(b, updater); err != nil {
		panic(err)
	}
	log.Printf("%s has been started...\n", b.Bot.User.Username)
	updater.Idle()
}
