package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func StartHandler() handlers.Command {
	return handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := ctx.EffectiveMessage.Reply(b, "Привет! Отправь любое сообщение, и я повторю его.", nil)
		return err
	})
}
