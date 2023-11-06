package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"syscall"

	"zpcg/internal/app"
)

const (
	TelegramApiTokenEnv = "TELEGRAM_APITOKEN"
	TimetableGobPathEnv = "TIMETABLE_GOB_PATH"
)

func main() {
	const logfmt = "main: "
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// config
	telegramApiToken, found := os.LookupEnv(TelegramApiTokenEnv)
	if !found {
		log.Fatal(logfmt, "can't find telegram api token env: ", TelegramApiTokenEnv)
	}
	timetableGobPath, found := os.LookupEnv(TimetableGobPathEnv)
	if !found {
		log.Fatal(logfmt, "can't find timetable gob path env env: ", TimetableGobPathEnv)
	}
	// tg bot
	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		log.Fatal(logfmt, "tgbotapi.NewBotAPI", err)
	}
	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)
	// app
	_app, err := app.NewApp(timetableGobPath)
	if err != nil {
		log.Fatal(err)
	}

	// start
	go receiveUpdates(ctx, _app, bot, updates)

	// wait for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	// done
	cancel()
}

func receiveUpdates(ctx context.Context, _app *app.App, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	const logfmt = "receiveUpdates: "
	for {
		if ctx.Err() != nil {
			return
		}
		msg, isNotEmpty := _app.HandleUpdate(<-updates)
		if !isNotEmpty {
			continue
		}
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(logfmt, "bot.Send: ", err)
		}
	}
}
