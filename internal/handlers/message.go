package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func EchoHandler() handlers.Message {
	return handlers.NewMessage(nil, func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
		return err
	})
}
