package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"

	"zpcg/internal/config"
	"zpcg/internal/model"
)

type App interface {
	HandleUpdate(update model.Update) (response []model.ResponseWithChatId, warning error)
}

// RunServer starts all processes needed to communicate with environment - initializes http server, logger,
// middlewares, k8s probes, etc. It knows nothing about business logic, only handles communication
func RunServer(ctx context.Context, _config config.Config, _app App) error {
	// logger
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)
	// tg bot
	bot, err := gotgbot.NewBot(_config.TelegramApiToken, &gotgbot.BotOpts{
		DisableTokenCheck: true,
	})
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

type PostTgMsgHandler func(...model.ResponseWithChatId) []model.ResponseWithChatId

func newUpdatesHandler(ctx context.Context, _app App, bot *gotgbot.Bot,
	messagePostHandler ...PostTgMsgHandler) func(http.ResponseWriter, *http.Request) {
	const logfmt = "receiveUpdates: "
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			update              gotgbot.Update
			responses           []string
			finalError, warning error
			httpStatus          = http.StatusOK
		)
		defer func() {
			logTrace(update, responses, warning, finalError)
			w.WriteHeader(httpStatus)
		}()

		// handle update
		if ctx.Err() != nil {
			finalError = ctx.Err()
			httpStatus = http.StatusInternalServerError
			return
		}
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			finalError = fmt.Errorf(logfmt+"json.NewDecoder(r.Body).Decode(&update): %w", err)
			httpStatus = http.StatusBadRequest
			return
		}
		var messages []model.ResponseWithChatId
		messages, warning = _app.HandleUpdate(UpdateFromTelegram(update))
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
			chatId, response, opts := ResponseToTelegram(finalMessage)
			responses = append(responses, response)
			_, err := bot.SendMessage(chatId, response, opts)
			if err != nil {
				sendError = errors.Join(sendError, err)
			}
		}
		if sendError != nil {
			finalError = fmt.Errorf(logfmt+"bot.SendMessage: %w", sendError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func logTrace(update gotgbot.Update, responseTexts []string, warning, finalError error) {
	const logFmt = "handleUpdate: %s"
	var (
		message  = update.Message
		logEvent *zerolog.Event
	)
	// set level
	switch {
	case finalError != nil:
		logEvent = log.Error().AnErr("error", finalError).AnErr("warning", warning)
	case warning != nil:
		logEvent = log.Warn().AnErr("warning", warning)
	default:
		logEvent = log.Trace()
	}
	// log responseTexts
	var responseShorts []string
	for _, response := range responseTexts {
		// get first 2 lines of the response
		responseLines := strings.Split(response, "\n")
		if len(responseLines) > 2 {
			responseLines = responseLines[:2]
		}
		responseShorts = append(responseShorts, strings.Join(responseLines, "\n"))
	}
	// log
	logEvent.
		Int64("chatId", message.Chat.Id).
		Str("languageCode", update.Message.From.LanguageCode).
		Str("messageText", message.Text).
		Strs("responseShorts", responseShorts).
		Msgf(logFmt, "new message handled")
}
