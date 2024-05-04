package main

import (
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	var token string
	if _token, ok := os.LookupEnv("TELEGRAM_APITOKEN"); ok {
		token = _token
	} else {
		log.Fatal("TELEGRAM_APITOKEN environment variable not set")
	}
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Fatal("gotgbot.NewBot: ", err)
	}

	bot.SetMyName()
	bot.SetMyDescription()
	bot.SetMyShortDescription()
	bot.SetMyCommands()

}
