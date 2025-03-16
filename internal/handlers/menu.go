package handlers

import (
	"context"
	"errors"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/application"
	"log"
)

var (
	ErrUserNotFound = errors.New("no user found")
)

func showStartMenu(b *gotgbot.Bot, ctx *ext.Context) {
	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{Text: "Профиль", CallbackData: "profile_callback"}},
			{{Text: "Настроить оповещения", CallbackData: "notifications_callback"}},
		},
	}
	_, _ = ctx.EffectiveMessage.Reply(b, "Выберите действие:", &gotgbot.SendMessageOpts{ReplyMarkup: keyboard})
}

func SettingsMenuCallbackHandler(db *application.UserService) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("profile_settings"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		user, err := db.GetUser(context.Background(), userID)
		if err != nil {
			log.Printf("error getting user: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка загрузки профиля.")
			return err
		}

		if user == nil {
			sendErrorMessage(b, cb, userID, "Ошибка загрузки профиля.")
			return ErrUserNotFound
		}

		keyboard := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Оповещения", CallbackData: "notifications_callback"}},
				{{Text: "Назад", CallbackData: "profile_callback"}},
			},
		}
		_, _, err = ctx.EffectiveMessage.EditText(b, "Выберите действие:", &gotgbot.EditMessageTextOpts{ReplyMarkup: keyboard})
		if err != nil {
			log.Printf("error editing message: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка загрузки меню настроек.")
			return err
		}
		return nil
	})
}
