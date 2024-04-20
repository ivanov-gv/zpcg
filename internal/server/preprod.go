package server

import (
	"zpcg/internal/model/message"
)

func AddTestEnvWarning(messages message.ResponseWithChatId) message.ResponseWithChatId {
	const warningText = "" +
		"THIS IS A TEST VERSION. DO NOT USE IT \n" +
		"THE RELEASE VERSION IS HERE: \n" +
		"\n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n"

	if len(messages.Send) == 0 {
		return messages
	}

	messages.Send = append(messages.Send, message.Response{
		Text: warningText,
	})
	return messages
}
