package server

import (
	"github.com/samber/lo"

	"zpcg/internal/model"
)

func AddTestEnvWarning(messages ...model.ResponseWithChatId) []model.ResponseWithChatId {
	const warningText = "" +
		"THIS IS A TEST VERSION. DO NOT USE IT \n" +
		"THE RELEASE VERSION IS HERE: \n" +
		"\n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n"

	chatIdSet := lo.SliceToMap(messages, func(item model.ResponseWithChatId) (int64, struct{}) {
		return item.ChatId, struct{}{}
	})

	for chatId := range chatIdSet {
		messages = append(messages, model.ResponseWithChatId{
			Response: model.Response{
				Text: warningText,
			},
			ChatId: chatId,
		})
	}
	return messages
}
