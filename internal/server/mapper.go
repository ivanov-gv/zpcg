package server

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/message"
)

func UpdateFromTelegram(update gotgbot.Update) message.Update {
	switch {
	case update.Message != nil:
		var _message = message.Message{
			IsFilled: true,
			Text:     update.Message.Text,
			ChatId:   update.Message.Chat.Id,
		}
		if update.Message.From != nil {
			_message.From = message.From{
				IsFilled:     true,
				LanguageCode: update.Message.From.LanguageCode,
			}
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

func ResponseToTelegramSend(response message.ToSend) (text string, opts *gotgbot.SendMessageOpts) {
	var replyMarkup gotgbot.ReplyMarkup
	if len(response.InlineKeyboard) != 0 {
		replyMarkup = inlineKeyboardToTelegram(response.InlineKeyboard)
	} else {
		// FIXME: removes keyboard for the users who has it from the older versions of the @monterails_bot
		//   we need to send RemoveKeyboard message silently to everyone with an update text after the brand new version will be implemented
		replyMarkup = gotgbot.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      false,
		}
	}

	return response.Text, &gotgbot.SendMessageOpts{
		ParseMode:   string(response.ParseMode),
		ReplyMarkup: replyMarkup,
	}
}

func ResponseToTelegramUpdate(chatId int64, response message.ToUpdate) (text string, opts *gotgbot.EditMessageTextOpts) {
	var (
		_chatId          int64
		_messageId       int64
		_inlineMessageId string
	)
	if len(response.InlineMessageId) == 0 {
		_chatId = chatId
		_messageId = response.MessageId
	} else {
		_inlineMessageId = response.InlineMessageId
	}

	text = response.Text
	opts = &gotgbot.EditMessageTextOpts{
		ChatId:          _chatId,
		MessageId:       _messageId,
		InlineMessageId: _inlineMessageId,
		ParseMode:       string(response.ParseMode),
		ReplyMarkup:     inlineKeyboardToTelegram(response.InlineKeyboard),
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

func ResponseToTelegramAnswerCallbackQuery(query message.ToAnswerCallbackQuery) (string, *gotgbot.AnswerCallbackQueryOpts) {
	var (
		callbackQueryId = query.CallbackQueryId
		opts            *gotgbot.AnswerCallbackQueryOpts
	)
	if len(query.Text) != 0 {
		opts = &gotgbot.AnswerCallbackQueryOpts{
			Text:      query.Text,
			ShowAlert: query.ShowAlert,
		}
	}
	return callbackQueryId, opts
}
