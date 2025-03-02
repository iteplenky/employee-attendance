package handlers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/application"
	"log"
)

func NotificationsSettingsHandler(db *application.UserService) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("notifications_callback"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		enabled, err := db.AreNotificationsEnabled(context.Background(), userID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error checking notifications enabled: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка получения статуса оповещений.")
			return err
		}

		_, _, err = cb.Message.EditText(b, "Настройки оповещений:", &gotgbot.EditMessageTextOpts{
			ReplyMarkup: getNotificationKeyboard(enabled),
		})
		if err != nil {
			log.Printf("error editing message: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка при изменении сообщения.")
		}
		return err
	})
}

func ToggleNotificationsHandler(db *application.UserService, cache *application.SubscriptionService) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("toggle_notifications"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		enabled, err := db.AreNotificationsEnabled(context.Background(), userID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error checking notifications enabled: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка получения статуса.")
			return err
		}

		newState := !enabled
		err = db.ToggleNotifications(context.Background(), userID, newState)
		if err != nil {
			log.Printf("error toggling notifications: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка обновления настроек.")
			return err
		}

		user, err := db.GetUser(context.Background(), userID)
		if err != nil {
			log.Printf("error getting user: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка получения пользователя.")
			return err
		}

		if user == nil {
			sendErrorMessage(b, cb, userID, "Ошибка получения пользователя.")
			return errors.New("no user found")
		}

		statusText := "Вы отписались от оповещений."
		if newState {
			statusText = "Вы подписаны на оповещения."
			err = cache.SaveSubscribersToCache(context.Background(), map[string]int64{user.IIN: userID})
		} else {
			err = cache.RemoveSubscriberFromCache(context.Background(), user.IIN)
		}
		if err != nil {
			log.Printf("error removing subscriber from cache: %v", err)
		}

		_, _, err = cb.Message.EditText(b, statusText, &gotgbot.EditMessageTextOpts{
			ReplyMarkup: getNotificationKeyboard(newState),
		})
		if err != nil {
			log.Printf("error updating notifications: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка при изменении сообщения.")
		}
		return err
	})
}

func getNotificationKeyboard(enabled bool) gotgbot.InlineKeyboardMarkup {
	buttonText := "🔔 Подписаться"
	if enabled {
		buttonText = "🔕 Отписаться"
	}

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: buttonText, CallbackData: "toggle_notifications"},
			},
			{
				{Text: "Назад в Настройки", CallbackData: "profile_settings"},
			},
		},
	}
}
