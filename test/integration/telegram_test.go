package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mock_server "github.com/ivanov-gv/zpcg/gen/mocks/server"
	"github.com/ivanov-gv/zpcg/internal/app"
	"github.com/ivanov-gv/zpcg/internal/config"
	model_render "github.com/ivanov-gv/zpcg/internal/model/render"
	"github.com/ivanov-gv/zpcg/internal/server"
)

const (
	TelegramApiToken  = "123:abcd"
	Port              = "8080" // TODO: take random open port
	Url               = "localhost:" + Port
	HttpServerAddress = "http://" + Url
	Environment       = "TestEnvironment"
)

var _config = config.Config{
	TelegramApiToken: TelegramApiToken,
	Port:             Port,
	Environment:      Environment,
}

// TestTelegramRouteResponse needs tdlib installed to run. see: https://tdlib.github.io/td/build.html
func TestTelegramRouteResponse(t *testing.T) {
	_app, err := app.NewApp()
	assert.NoError(t, err)
	mockTgClient := &mock_server.MockCustomTgClient{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start server
	go func() {
		err := server.RunServer(ctx, _config, _app, server.WithCustomTgClient(mockTgClient))
		assert.NotErrorIs(t, err, http.ErrServerClosed)
	}()
	assert.Eventually(t, func() bool {
		conn, err := net.Dial("tcp", Url)
		if conn != nil {
			defer conn.Close()
		}
		return err == nil
	}, time.Second, 100*time.Millisecond)

	// run tests
	// send empty message update
	t.Run("message update", func(t *testing.T) {
		const request = "{}"
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer([]byte(request)))
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, response.StatusCode, http.StatusOK)
	})
	// start message
	t.Run("start message", func(t *testing.T) {
		for _, languageTag := range model_render.SupportedLanguages {
			language := languageTag.String()
			t.Run(language, func(t *testing.T) {
				request := gotgbot.Update{
					Message: &gotgbot.Message{
						Text: "/start",
						From: &gotgbot.User{
							LanguageCode: language,
						},
					},
				}
				requestRaw := lo.Must(json.Marshal(request))
				mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "sendMessage", mock.Anything, mock.Anything, mock.Anything).
					Return([]byte("{}"), nil)
				response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
				assert.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, response.StatusCode, http.StatusOK)
				t.Log("telegram response: ", lo.LastOrEmpty(mockTgClient.Calls).Arguments)
			})
		}
	})
	// timetable request
	t.Run("timetable request", func(t *testing.T) {
		request := gotgbot.Update{
			Message: &gotgbot.Message{
				Text: "nikshich, bar",
				From: &gotgbot.User{
					LanguageCode: "en",
				},
			},
		}
		requestRaw := lo.Must(json.Marshal(request))
		//mockTgClient.EXPECT().TimeoutContext(mock.Anything).Return(context.WithCancel(ctx))
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "sendMessage", mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("{}"), nil)
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, response.StatusCode, http.StatusOK)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
	// callback update
	t.Run("update callback", func(t *testing.T) {
		request := gotgbot.Update{
			CallbackQuery: &gotgbot.CallbackQuery{
				Id:      "someid",
				Data:    "0|Niksic|Bar|2000-01-01",
				Message: &gotgbot.Message{},
			},
		}

		requestRaw := lo.Must(json.Marshal(request))
		//mockTgClient.EXPECT().TimeoutContext(mock.Anything).Return(context.WithCancel(ctx))
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "answerCallbackQuery", mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("true"), nil)
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "editMessageText", mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("{}"), nil)
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
	// callback reverse
	t.Run("reverse route callback", func(t *testing.T) {
		request := gotgbot.Update{
			CallbackQuery: &gotgbot.CallbackQuery{
				Id:      "someid",
				Data:    "1|Niksic|Bar",
				Message: &gotgbot.Message{},
			},
		}

		requestRaw := lo.Must(json.Marshal(request))
		//mockTgClient.EXPECT().TimeoutContext(mock.Anything).Return(context.WithCancel(ctx))
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "answerCallbackQuery", mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("true"), nil)
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "sendMessage", mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("{}"), nil)
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
}
