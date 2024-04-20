package server

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/samber/lo"

	"zpcg/internal/model/message"
)

func UpdateFromTelegram(update gotgbot.Update) message.Update {
	switch {
	case update.Message != nil:
		var _message message.Message
		if update.Message.From != nil {
			_message.From = message.From{
				IsFilled:     true,
				LanguageCode: update.Message.From.LanguageCode,
			}
		}
		_message = message.Message{
			IsFilled: true,
			Text:     update.Message.Text,
			ChatId:   update.Message.Chat.Id,
		}
		return message.Update{
			Type:    message.MessageUpdateType,
			Message: _message,
		}
	case update.CallbackQuery != nil:
		callback := update.CallbackQuery
		return message.Update{
			Type: message.CallbackUpdateType,
			Callback: message.Callback{
				Id:     callback.Id,
				ChatId: callback.Message.GetChat().Id,
				From: message.From{
					IsFilled:     true,
					LanguageCode: callback.From.LanguageCode,
				},
				Message: message.MaybeInaccessibleMessage{
					Id: callback.Message.GetMessageId(),
				},
				InlineMessageId: callback.InlineMessageId,
				Data:            callback.Data,
			},
		}
	default:
		return message.Update{Type: message.UnsupportedUpdateType}
	}
}

func ResponseToTelegramSend(response message.Response) (text string, opts *gotgbot.SendMessageOpts) {
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
	return text, opts
}

func ResponseToTelegramUpdate(chatId int64, response message.ToUpdate) (text string, opts *gotgbot.EditMessageTextOpts) {
	text = response.Text
	opts = &gotgbot.EditMessageTextOpts{
		ChatId:    chatId,
		MessageId: response.MessageId,
		// TODO: inline message id might also be used. is it a better approach?
		ParseMode:   string(response.ParseMode),
		ReplyMarkup: inlineKeyboardToTelegram(response.InlineKeyboard),
	}

	return text, opts
}

func inlineKeyboardToTelegram(inlineKeyboard [][]message.InlineButton) gotgbot.InlineKeyboardMarkup {
	var inlineKeyboardOutput [][]gotgbot.InlineKeyboardButton
	for _, row := range inlineKeyboard {
		var inlineRow []gotgbot.InlineKeyboardButton
		for _, button := range row {
			inlineRow = append(inlineRow, lo.Must(inlineButtonToTelegram(button)))
		}
		inlineKeyboardOutput = append(inlineKeyboardOutput, inlineRow)
	}
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboardOutput,
	}
}

func inlineButtonToTelegram(button message.InlineButton) (gotgbot.InlineKeyboardButton, error) {
	switch button.Type {
	case message.UrlInlineButtonType:
		return gotgbot.InlineKeyboardButton{
			Text: button.Text,
			Url:  button.Url.Url,
		}, nil
	case message.CallbackInlineButtonType:
		return gotgbot.InlineKeyboardButton{
			Text:         button.Text,
			CallbackData: button.Callback.Data,
		}, nil
	default:
		return gotgbot.InlineKeyboardButton{}, fmt.Errorf("unsupported button type: %v", button.Type)
	}
}
