package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"

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
		log.Fatal().Err(fmt.Errorf(logfmt+"can't find telegram api token env: %s", TelegramApiTokenEnv))
	}
	serverPort, found := os.LookupEnv(PortEnv)
	if !found {
		log.Fatal().Err(fmt.Errorf(logfmt+"can't find server port env env: %s", PortEnv))
	}
	timetableReader, err := resources.FS.Open(TimetableGobFileName)
	if err != nil {
		log.Fatal().Err(fmt.Errorf(logfmt+"fs.Open: %w", err))
	}
	// tg bot
	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		log.Fatal().Err(fmt.Errorf(logfmt+"tgbotapi.NewBotAPI: %w", err))
	}
	// app
	_app, err := app.NewApp(timetableReader)
	if err != nil {
		log.Fatal().Err(err)
	}
	// logger
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)
	// server
	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(newUpdatesHandler(ctx, _app, bot))))
	mux.HandleFunc("/health", func(_ http.ResponseWriter, _ *http.Request) { return })
	log.Printf("listening on port %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err)
	}
}

func newUpdatesHandler(ctx context.Context, _app *app.App, bot *tgbotapi.BotAPI) func(http.ResponseWriter, *http.Request) {
	const logfmt = "receiveUpdates: "
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.Err() != nil {
			return
		}
		var (
			update tgbotapi.Update
			logger = log.Ctx(r.Context())
		)
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			logger.Error().Err(fmt.Errorf(logfmt+"json.NewDecoder(r.Body).Decode(&update): %w", err)).Send()
			return
		}
		msg, isNotEmpty := _app.HandleUpdate(update)
		if !isNotEmpty {
			return
		}
		_, err = bot.Send(msg)
		if err != nil {
			logger.Error().Err(fmt.Errorf(logfmt+"bot.Send: %w", err)).Send()
			return
		}
	}
}
