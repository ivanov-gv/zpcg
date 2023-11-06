package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"zpcg/internal/app"
	"zpcg/resources"
)

const (
	TelegramApiTokenEnv  = "TELEGRAM_APITOKEN"
	PortEnv              = "PORT"
	TimetableGobFileName = "timetable.gob"
)

func main() {
	const logfmt = "main: "
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// config
	telegramApiToken, found := os.LookupEnv(TelegramApiTokenEnv)
	if !found {
		log.Fatal(logfmt, "can't find telegram api token env: ", TelegramApiTokenEnv)
	}
	serverPort, found := os.LookupEnv(PortEnv)
	if !found {
		log.Fatal(logfmt, "can't find server port env env: ", PortEnv)
	}
	timetableReader, err := resources.FS.Open(TimetableGobFileName)
	if err != nil {
		log.Fatal(logfmt, "fs.Open", err)
	}
	// tg bot
	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		log.Fatal(logfmt, "tgbotapi.NewBotAPI ", err)
	}
	// app
	_app, err := app.NewApp(timetableReader)
	if err != nil {
		log.Fatal(err)
	}
	// server
	mux := http.NewServeMux()
	mux.HandleFunc("/", newUpdatesHandler(ctx, _app, bot))
	mux.HandleFunc("/health", func(_ http.ResponseWriter, _ *http.Request) { return })
	log.Printf("listening on port %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func newUpdatesHandler(ctx context.Context, _app *app.App, bot *tgbotapi.BotAPI) func(http.ResponseWriter, *http.Request) {
	const logfmt = "receiveUpdates: "
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.Err() != nil {
			return
		}
		var update tgbotapi.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Println(logfmt, "json.NewDecoder(r.Body).Decode(&update): ", err)
			return
		}
		msg, isNotEmpty := _app.HandleUpdate(update)
		if !isNotEmpty {
			return
		}
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(logfmt, "bot.Send: ", err)
			return
		}
	}
}
