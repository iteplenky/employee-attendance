package handlers

import (
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/iteplenky/employee-attendance/application"
	"log"
	"sync"
	"unicode/utf8"
)

var userStates sync.Map

func StartHandler(db *application.UserService) handlers.Command {
	return handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := ctx.EffectiveUser.Id

		if _, exists := userStates.Load(userID); exists {
			_, _ = ctx.EffectiveMessage.Reply(b, "Вы уже в процессе регистрации. Введите ваш ИИН:", nil)
			return nil
		}

		user, err := db.GetUser(context.Background(), userID)
		if err != nil {
			log.Printf("error getting user: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка при получении данных пользователя.")
			return err
		}

		if user != nil {
			showStartMenu(b, ctx)
		} else {
			userStates.Store(userID, true)
			_, _ = ctx.EffectiveMessage.Reply(b, "Вы не зарегистрированы. Введите ваш ИИН:", nil)
		}
		return nil
	})
}

func IINHandler(db *application.UserService) handlers.Message {
	return handlers.NewMessage(nil, func(b *gotgbot.Bot, ctx *ext.Context) error {
		cb := ctx.Update.CallbackQuery
		userID := ctx.EffectiveUser.Id

		if _, exists := userStates.Load(userID); !exists {
			return nil
		}

		iin := ctx.EffectiveMessage.Text
		iinLen := utf8.RuneCountInString(iin)

		if iinLen < 10 || iinLen > 14 {
			_, _ = ctx.EffectiveMessage.Reply(b, "Некорректная длина ИИН, введите еще раз.", nil)
			return nil
		}

		err := db.RegisterUser(context.Background(), userID, iin)
		if err != nil {
			log.Printf("error registering user: %v\n", err)
			sendErrorMessage(b, cb, userID, "Ошибка при регистрации, попробуйте позже.")
			return err
		}

		userStates.Delete(userID)
		showStartMenu(b, ctx)
		return nil
	})
}
