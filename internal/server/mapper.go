package server

import (
	"github.com/PaulSonOfLars/gotgbot/v2"

	"zpcg/internal/model"
)

func UpdateFromTelegram(update gotgbot.Update) model.Update {
	var (
		message model.Message
		from    model.From
	)
	if update.Message != nil && update.Message.From != nil {
		from = model.From{
			IsFilled:     true,
			LanguageCode: update.Message.From.LanguageCode,
		}
	}
	if update.Message != nil {
		message = model.Message{
			IsFilled: true,
			From:     from,
			Text:     update.Message.Text,
			ChatId:   update.Message.Chat.Id,
		}
	}
	return model.Update{
		Message: message,
	}
}

func ResponseToTelegram(response model.ResponseWithChatId) (chatId int64, text string, opts *gotgbot.SendMessageOpts) {
	chatId = response.ChatId
	text = response.Text
	opts = &gotgbot.SendMessageOpts{
		ParseMode:   string(response.ParseMode),
		ReplyMarkup: inlineKeyboardToTelegram(response.InlineKeyboard),
	}

	// FIXME: removes keyboard for the users who has it from the older versions of the @monterails_bot
	//   we need to send RemoveKeyboard message silently to everyone with an update text after the brand new version will be implemented
	if opts.ReplyMarkup == nil {
		opts.ReplyMarkup = gotgbot.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      false,
		}
	}
	return chatId, text, opts
}

func inlineKeyboardToTelegram(inlineKeyboard [][]model.InlineButton) gotgbot.ReplyMarkup {
	if inlineKeyboard == nil {
		return nil
	}
	var inlineKeyboardOutput [][]gotgbot.InlineKeyboardButton
	for _, row := range inlineKeyboard {
		var inlineRow []gotgbot.InlineKeyboardButton
		for _, button := range row {
			inlineRow = append(inlineRow,
				gotgbot.InlineKeyboardButton{
					Text: button.Text,
					Url:  button.Url,
				})
		}
		inlineKeyboardOutput = append(inlineKeyboardOutput, inlineRow)
	}
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboardOutput,
	}
}
