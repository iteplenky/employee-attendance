package handlers

import "github.com/PaulSonOfLars/gotgbot/v2"

func sendErrorMessage(b *gotgbot.Bot, cb *gotgbot.CallbackQuery, userID int64, text string) {
	_, _ = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: text})
	if cb.Message != nil {
		_, _ = b.DeleteMessages(userID, []int64{cb.Message.GetMessageId()}, nil)
	}
}
