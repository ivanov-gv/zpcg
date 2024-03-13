package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"zpcg/internal/model"
)

func ResponseToTelegram(chatId int64, response model.Response) tgbotapi.MessageConfig {
	// create message and return
	output := tgbotapi.NewMessage(chatId, response.Text)
	output.ParseMode = response.ParseMode
	output.ReplyMarkup = inlineKeyboardToTelegram(response.InlineKeyboard)

	// FIXME: removes keyboard for the users who has it from the older versions of the @monterails_bot
	//   we need to send RemoveKeyboard message silently to everyone with an update text after the brand new version will be implemented
	if output.ReplyMarkup != nil {
		output.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}
	return output
}

func inlineKeyboardToTelegram(inlineKeyboard [][]model.InlineButton) any {
	if inlineKeyboard == nil {
		return nil
	}
	var inlineKeyboardOutput [][]tgbotapi.InlineKeyboardButton
	for _, row := range inlineKeyboard {
		var inlineRow []tgbotapi.InlineKeyboardButton
		for _, button := range row {
			inlineRow = append(inlineRow, tgbotapi.NewInlineKeyboardButtonURL(button.Text, button.Url))
		}
		inlineKeyboardOutput = append(inlineKeyboardOutput, inlineRow)
	}
	return tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardOutput...)
}
