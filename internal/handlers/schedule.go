package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/database"
)

func ScheduleCallbackHandler() handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("schedule_callback"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery

		keyboard := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: "09:00 - 18:00", CallbackData: "schedule_0900_1800"},
					{Text: "10:00 - 19:00", CallbackData: "schedule_1000_1900"},
				},
				{
					{Text: "Профиль", CallbackData: "profile_callback"},
				},
			},
		}

		_, _, err := cb.Message.EditText(b, "Для того чтобы получать оповещения, необходимо выбрать рабочий график:", &gotgbot.EditMessageTextOpts{
			ReplyMarkup: keyboard,
		})
		return err
	})
}

func ScheduleSelectionHandler(db database.UserRepository) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Prefix("schedule_"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		schedule := cb.Data[len("schedule_"):]

		start := schedule[:2] + ":" + schedule[2:4]
		end := schedule[5:7] + ":" + schedule[7:9]

		err := db.SaveSchedule(cb.From.Id, start, end)
		if err != nil {
			_, _ = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Ошибка сохранения графика!"})
			return err
		}

		msg := "Вы выбрали график: " + start + " - " + end
		_, _, err = cb.Message.EditText(b, msg, &gotgbot.EditMessageTextOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Профиль", CallbackData: "profile_callback"}},
			}}})
		return err
	})
}
