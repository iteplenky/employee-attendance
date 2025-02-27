package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/database"
)

func NotificationsSettingsHandler(db database.UserRepository) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("notifications_callback"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		enabled, err := db.AreNotificationsEnabled(userID)
		if err != nil {
			_, _ = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Ошибка получения статуса оповещений."})
			return err
		}

		buttonText := "Подписаться"
		if enabled {
			buttonText = "Отписаться"
		}

		keyboard := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: buttonText, CallbackData: "toggle_notifications"},
				},
				{
					{Text: "Профиль", CallbackData: "profile_callback"},
				},
			},
		}

		_, _, err = cb.Message.EditText(b, "Настройки оповещений:", &gotgbot.EditMessageTextOpts{
			ReplyMarkup: keyboard,
		})
		return err
	})
}

func ToggleNotificationsHandler(db database.UserRepository) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("toggle_notifications"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		enabled, err := db.AreNotificationsEnabled(userID)
		if err != nil {
			_, _ = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Ошибка получения статуса."})
			return err
		}

		newState := !enabled
		err = db.ToggleNotifications(userID, newState)
		if err != nil {
			_, _ = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Ошибка обновления настроек."})
			return err
		}

		buttonText := "Подписаться"
		statusText := "Вы отписались от оповещений."
		if newState {
			buttonText = "Отписаться"
			statusText = "Вы подписаны на оповещения."
		}

		keyboard := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: buttonText, CallbackData: "toggle_notifications"},
				},
				{
					{Text: "Профиль", CallbackData: "profile_callback"},
				},
			},
		}

		_, _, err = cb.Message.EditText(b, statusText, &gotgbot.EditMessageTextOpts{
			ReplyMarkup: keyboard,
		})
		return err
	})
}
