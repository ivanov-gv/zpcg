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
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"

	"github.com/ivanov-gv/zpcg/internal/config"
	"github.com/ivanov-gv/zpcg/internal/model/message"
)

type App interface {
	HandleUpdate(update message.Update) (response message.ResponseWithChatId, warning error)
}

// RunServer starts all processes needed to communicate with environment - initializes http server, logger,
// middlewares, k8s probes, etc. It knows nothing about business logic, only handles communication
func RunServer(ctx context.Context, _config config.Config, _app App, opts ...ApplyOption) error {
	// settings options
	var settings = options{
		botOpts: &gotgbot.BotOpts{
			DisableTokenCheck: true,
		},
	}
	settings = applyOptions(settings, opts...)
	// logger
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)
	// tg bot
	bot, err := gotgbot.NewBot(_config.TelegramApiToken, settings.botOpts)
	if err != nil {
		return fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}
	// updates cache
	updateCache, err := lru.New[int64, int8](64)
	if err != nil {
		return fmt.Errorf("lru.New[int64, int8]: %w", err)
	}

	// server middlewares
	var postHandlers []PostTgMsgHandler
	if _config.Environment == config.EnvironmentPreProdValue {
		// add test env warning to every message
		postHandlers = append(postHandlers, AddTestEnvWarning)
	}
	// server
	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(newUpdatesHandler(ctx, _app, bot, updateCache, postHandlers...))))
	mux.HandleFunc("/health", func(_ http.ResponseWriter, _ *http.Request) { return })
	// start
	if err := http.ListenAndServe(":"+_config.Port, mux); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}

type PostTgMsgHandler func(message.ResponseWithChatId) message.ResponseWithChatId

const (
	maxUpdateRetry = 3
)

func newUpdatesHandler(ctx context.Context, _app App, bot *gotgbot.Bot, updateCache *lru.Cache[int64, int8],
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
		// retry updates only certain number of times
		value, ok := updateCache.Get(update.UpdateId)
		if ok && value >= maxUpdateRetry {
			return
		}
		updateCache.Add(update.UpdateId, value+1)

		var messages message.ResponseWithChatId
		messages, warning = _app.HandleUpdate(UpdateFromTelegram(update))
		if len(messages.Update) == 0 && len(messages.Send) == 0 && len(messages.Delete) == 0 && len(messages.AnswerCallback.CallbackQueryId) == 0 {
			return
		}
		// post update handlers
		for _, handler := range messagePostHandler {
			if handler == nil {
				continue
			}
			messages = handler(messages)
		}
		// send
		for _, sendMessage := range messages.Send {
			response, opts := ResponseToTelegramSend(sendMessage)
			responses = append(responses, response)
			_, err = bot.SendMessage(messages.ChatId, response, opts)
			if err != nil {
				finalError = fmt.Errorf(logfmt+"bot.SendMessage: %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// update
		for _, updateMessage := range messages.Update {
			response, opts := ResponseToTelegramUpdate(messages.ChatId, updateMessage)
			responses = append(responses, response)
			_, _, err = bot.EditMessageText(response, opts)
			if err != nil {
				finalError = fmt.Errorf(logfmt+"bot.EditMessageText: %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// answer callback
		if len(messages.AnswerCallback.CallbackQueryId) != 0 {
			_, err = bot.AnswerCallbackQuery(ResponseToTelegramAnswerCallbackQuery(messages.AnswerCallback))
			if err != nil {
				finalError = fmt.Errorf(logfmt+"bot.AnswerCallbackQuery: %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func logTrace(update gotgbot.Update, responseTexts []string, warning, finalError error) {
	const logFmt = "handleUpdate: %s"
	var logEvent *zerolog.Event
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
	switch {
	case update.Message != nil:
		chatId := update.Message.Chat.Id / 10000 // cut last 4 digits
		languageCode := update.Message.From.LanguageCode
		text := update.Message.Text
		// log
		logEvent.
			Int64("chatId", chatId).
			Str("languageCode", languageCode).
			Str("messageText", text).
			Strs("responseShorts", responseShorts).
			Msgf(logFmt, "new message handled")
	case update.CallbackQuery != nil:
		chatId := update.CallbackQuery.Message.GetChat().Id / 10000 // cut last 4 digits
		languageCode := update.CallbackQuery.From.LanguageCode
		callbackData := update.CallbackQuery.Data

		var text string
		if _message, ok := update.CallbackQuery.Message.(gotgbot.Message); ok {
			text = _message.Text
		}

		// log
		logEvent.
			Int64("chatId", chatId).
			Str("languageCode", languageCode).
			Str("messageText", text).
			Str("callbackData", callbackData).
			Strs("responseShorts", responseShorts).
			Msgf(logFmt, "new callback handled")
	}
}
