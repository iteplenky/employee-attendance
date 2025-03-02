package handlers

import (
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/iteplenky/employee-attendance/application"
	"log"
	"time"
)

func AttendanceCallbackHandler(db *application.UserService, fetcher *application.FetcherService) handlers.CallbackQuery {
	return handlers.NewCallback(callbackquery.Equal("attendance_callback"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := cb.From.Id

		user, err := db.GetUser(context.Background(), userID)
		if err != nil || user == nil {
			log.Printf("error getting user: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка загрузки пользователя.")
			return err
		}

		attendances, err := fetcher.GetUserAttendanceRecords(context.Background(), user.IIN)
		if err != nil {
			log.Printf("error fetching attendance: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка получения явок.")
			return err
		}

		keyboard := &gotgbot.EditMessageTextOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Назад", CallbackData: "profile_callback"}},
			},
		}}

		if len(attendances) == 0 {
			_, _, err = cb.Message.EditText(b, "Ваши явки за текущий месяц отсутствуют.", keyboard)
			if err != nil {
				log.Printf("error updating message: %v\n", keyboard)
				sendErrorMessage(b, cb, userID, "Ошибка при обновлении сообщения.")
			}
			return nil
		}

		groupedAttendances := make(map[string][]string)
		for _, att := range attendances {
			date, timeStr := splitDateTime(att.PunchTime)
			groupedAttendances[date] = append(groupedAttendances[date], timeStr)
		}

		msg := "Ваши явки за текущий месяц:\n"
		for date, times := range groupedAttendances {
			msg += fmt.Sprintf("Дата: %s\nВремя: %s\n\n", date, fmt.Sprintf("%s", times))
		}

		_, _, err = cb.Message.EditText(b, msg, keyboard)
		if err != nil {
			log.Printf("error updating message: %v\n", keyboard)
			sendErrorMessage(b, cb, userID, "Ошибка при обновлении сообщения.")
		}
		return err
	})
}

func splitDateTime(punchTime string) (string, string) {
	t, err := time.Parse("2006-01-02T15:04:05.999999Z", punchTime)
	if err != nil {
		log.Printf("failed to parse time: %v", err)
		return punchTime, ""
	}
	return t.Format("02.01"), t.Format("15:04:05")
}
