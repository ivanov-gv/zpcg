package server

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

	"zpcg/internal/config"
)

type App interface {
	HandleUpdate(update tgbotapi.Update) []tgbotapi.MessageConfig
}

func RunServer(ctx context.Context, _config config.Config, _app App) error {
	// logger
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)
	// tg bot
	bot, err := tgbotapi.NewBotAPI(_config.TelegramApiToken)
	if err != nil {
		return fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}

	// server middlewares
	var postHandlers []PostTgMsgHandler
	if _config.Environment == config.EnvironmentPreProdValue {
		// add test env warning to every message
		postHandlers = append(postHandlers, AddTestEnvWarning)
	}
	// server
	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(newUpdatesHandler(ctx, _app, bot, postHandlers...))))
	mux.HandleFunc("/health", func(_ http.ResponseWriter, _ *http.Request) { return })
	// start
	if err := http.ListenAndServe(":"+_config.Port, mux); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}

type PostTgMsgHandler func(...tgbotapi.MessageConfig) []tgbotapi.MessageConfig

func newUpdatesHandler(ctx context.Context, _app App, bot *tgbotapi.BotAPI,
	messagePostHandler ...PostTgMsgHandler) func(http.ResponseWriter, *http.Request) {
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
		messages := _app.HandleUpdate(update)
		if len(messages) == 0 {
			return
		}
		// post update handlers
		for _, handler := range messagePostHandler {
			if handler == nil {
				continue
			}
			messages = handler(messages...)
		}
		// send
		var sendError error
		for _, finalMessage := range messages {
			_, err := bot.Send(finalMessage)
			if err != nil {
				sendError = errors.Join(sendError, err)
			}
		}
		if sendError != nil {
			logger.Error().Err(fmt.Errorf(logfmt+"bot.Send: %w", err)).Send()
			return
		}
	}
}
