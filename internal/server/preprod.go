package server

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func AddTestEnvWarning(messages ...tgbotapi.MessageConfig) []tgbotapi.MessageConfig {
	const warningText = "" +
		"THIS IS A TEST VERSION. DO NOT USE IT \n" +
		"FULL VERSION IS HERE: \n" +
		"\n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n"

	chatIdSet := lo.SliceToMap(messages, func(item tgbotapi.MessageConfig) (int64, struct{}) {
		return item.ChatID, struct{}{}
	})

	for chatId := range chatIdSet {
		messages = append(messages, tgbotapi.NewMessage(chatId, warningText))
	}
	return messages
}
