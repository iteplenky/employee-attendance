package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/iteplenky/employee-attendance/database"
	"log"
	"sync"
	"unicode/utf8"
)

var userStates sync.Map

func StartHandler(db database.UserRepository) handlers.Command {
	return handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		userID := ctx.EffectiveUser.Id

		if _, exists := userStates.Load(userID); exists {
			_, _ = ctx.EffectiveMessage.Reply(b, "Вы уже в процессе регистрации. Введите ваш ИИН:", nil)
			return nil
		}

		exists, err := db.UserExists(userID)
		if err != nil {
			_, _ = ctx.EffectiveMessage.Reply(b, "Попробуйте позже.", nil)
			log.Printf("error checking if user exists: %v", err)
			return err
		}

		if exists {
			_, _ = ctx.EffectiveMessage.Reply(b, "Привет! Вы уже зарегистрированы.", nil)
		} else {
			userStates.Store(userID, true)
			_, _ = ctx.EffectiveMessage.Reply(b, "Вы не зарегистрированы. Пожалуйста, введите ваш ИИН:", nil)
		}
		return nil
	})
}

func IINHandler(db database.UserRepository) handlers.Message {
	return handlers.NewMessage(nil, func(b *gotgbot.Bot, ctx *ext.Context) error {
		userID := ctx.EffectiveUser.Id

		if _, exists := userStates.Load(userID); !exists {
			return nil
		}

		iin := ctx.EffectiveMessage.Text
		iinLen := utf8.RuneCountInString(iin)

		if iinLen < 8 {
			_, _ = ctx.EffectiveMessage.Reply(b, "Недостаточная длина ИИН, перепроверьте и введите еще раз.", nil)
			return nil
		} else if iinLen > 12 {
			_, _ = ctx.EffectiveMessage.Reply(b, "Большая длина ИИН, перепроверьте и введите еще раз.", nil)
			return nil
		}

		err := db.RegisterUser(userID, iin)
		if err != nil {
			_, _ = ctx.EffectiveMessage.Reply(b, "Попробуйте позже.", nil)
			log.Printf("error registering user: %v", err)
			return err
		}

		userStates.Delete(userID)
		_, _ = ctx.EffectiveMessage.Reply(b, "Регистрация успешна! Добро пожаловать!", nil)
		log.Printf("Registered user: %v", userID)
		return nil
	})
}
