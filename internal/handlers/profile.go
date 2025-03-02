package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/application"
	"log"
)

func ProfileCallbackHandler(db *application.UserService) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("profile_callback"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := ctx.Update.CallbackQuery.From.Id

		user, err := db.GetUser(context.Background(), userID)
		if err != nil {
			log.Printf("error getting user: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка загрузки профиля.")
			return err
		}

		if user == nil {
			sendErrorMessage(b, cb, userID, "Ошибка загрузки профиля.")
			return errors.New("no user found")
		}

		notificationText := "Не подписаны"
		if user.NotificationsEnabled {
			notificationText = "Подписаны"
		}

		msg := fmt.Sprintf("Ваш профиль:\nИИН: %s\nОповещения: %s", user.IIN, notificationText)
		_, _, err = ctx.Update.CallbackQuery.Message.EditText(b, msg, &gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{{Text: "Мои явки", CallbackData: "attendance_callback"}},
					{{Text: "Настройки", CallbackData: "profile_settings"}},
				},
			},
		})
		if err != nil {
			log.Printf("error updating message: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка при изменении сообщения.")
		}
		return err
	})
}
