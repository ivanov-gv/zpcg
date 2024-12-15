package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/render"
)

func main() {
	var token string
	if _token, ok := os.LookupEnv("TELEGRAM_APITOKEN"); ok {
		token = _token
	} else {
		log.Fatal("TELEGRAM_APITOKEN environment variable not set")
	}
	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: gotgbot.BotClient(&gotgbot.BaseBotClient{
			Client:             http.Client{},
			UseTestEnvironment: false,
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 30,
				APIURL:  "",
			},
		}),
		DisableTokenCheck: false,
		RequestOpts: &gotgbot.RequestOpts{
			Timeout: time.Second * 30,
			APIURL:  "",
		},
	})
	if err != nil {
		log.Fatal("gotgbot.NewBot: ", err)
	}

	// set info

	defaultInfo := render.BotInfoMap[render.DefaultLanguageTag]
	_, err = bot.SetMyName(&gotgbot.SetMyNameOpts{Name: defaultInfo.Name})
	if err != nil {
		log.Fatal("bot.SetMyName: ", err)
	}
	_, err = bot.SetMyDescription(&gotgbot.SetMyDescriptionOpts{Description: defaultInfo.Description})
	if err != nil {
		log.Fatal("bot.SetMyName: ", err)
	}
	_, err = bot.SetMyShortDescription(&gotgbot.SetMyShortDescriptionOpts{ShortDescription: defaultInfo.ShortDescription})
	if err != nil {
		log.Fatal("bot.SetMyShortDescription: ", err)
	}
	_, err = bot.SetMyCommands(lo.Map(render.AllCommands, func(item render.BotCommand, _ int) gotgbot.BotCommand {
		return gotgbot.BotCommand{
			Command:     string(item),
			Description: defaultInfo.CommandNames[item],
		}
	}), &gotgbot.SetMyCommandsOpts{Scope: gotgbot.BotCommandScopeDefault{}})
	if err != nil {
		log.Fatal("bot.SetMyCommands: ", err)
	}

	for languageTag, info := range render.BotInfoMap {
		_, err = bot.SetMyName(&gotgbot.SetMyNameOpts{
			Name:         info.Name,
			LanguageCode: languageTag.String(),
		})
		if err != nil {
			log.Fatal(languageTag.String()+" bot.SetMyName: ", err)
		}
		_, err = bot.SetMyDescription(&gotgbot.SetMyDescriptionOpts{
			Description:  info.Description,
			LanguageCode: languageTag.String(),
		})
		if err != nil {
			log.Fatal(languageTag.String()+" bot.SetMyName: ", err)
		}
		_, err = bot.SetMyShortDescription(&gotgbot.SetMyShortDescriptionOpts{
			ShortDescription: info.ShortDescription,
			LanguageCode:     languageTag.String(),
		})
		if err != nil {
			log.Fatal(languageTag.String()+" bot.SetMyShortDescription: ", err)
		}
		_, err = bot.SetMyCommands(lo.Map(render.AllCommands, func(item render.BotCommand, _ int) gotgbot.BotCommand {
			return gotgbot.BotCommand{
				Command:     string(item),
				Description: info.CommandNames[item],
			}
		}), &gotgbot.SetMyCommandsOpts{Scope: gotgbot.BotCommandScopeDefault{}, LanguageCode: languageTag.String()})
		if err != nil {
			log.Fatal(languageTag.String()+" bot.SetMyCommands: ", err)
		}
	}

}
