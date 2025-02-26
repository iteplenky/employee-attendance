package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/database"
)

func ProfileCallbackHandler(db database.UserRepository) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("profile_callback"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		user, err := db.GetUser(userID)
		if err != nil {
			_, _ = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Ошибка получения профиля."})
			return err
		}

		msg := fmt.Sprintf("Ваш профиль:\nИИН: %s", user.IIN)
		_, _, err = cb.Message.EditText(b, msg, nil)
		if err != nil {
			return err
		}

		return nil
	})
}
