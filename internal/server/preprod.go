package server

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddTestEnvWarning(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	const format = "" +
		"%s \n " +
		"\n" +
		"THIS IS A TEST VERSION. DO NOT USE IT \n" +
		"ЭТО ТЕСТОВАЯ ВЕРСИЯ. НЕ ИСПОЛЬЗУЙТЕ ЕЕ \n" +
		"\n" +
		"FULL VERSION IS HERE: @Monterails_bot \n" +
		"ПОЛНАЯ ВЕРСИЯ ЗДЕСЬ: @Monterails_bot \n" +
		"\n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n" +
		"@Monterails_bot \n"
	message.Text = fmt.Sprintf(format, message.Text)
	return message
}
