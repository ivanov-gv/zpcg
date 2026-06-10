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
	"golang.org/x/text/language"

	mock_server "github.com/ivanov-gv/zpcg/gen/mocks/server"
	timetable_gen "github.com/ivanov-gv/zpcg/gen/timetable"
	"github.com/ivanov-gv/zpcg/internal/app"
	"github.com/ivanov-gv/zpcg/internal/config/server_config"
	model_render "github.com/ivanov-gv/zpcg/internal/model/message_render"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/server"
	"github.com/ivanov-gv/zpcg/internal/service/date"
)

const (
	TelegramApiToken  = "123:abcd"
	Port              = "8080" // TODO: take random open port
	Url               = "localhost:" + Port
	HttpServerAddress = "http://" + Url
	Environment       = "TestEnvironment"
)

var (
	_config = server_config.Config{
		TelegramApiToken: TelegramApiToken,
		Port:             Port,
		Environment:      Environment,
	}
	updateId int64 = 0
)

func nextUpdateRequest() gotgbot.Update {
	updateId += 1
	return gotgbot.Update{UpdateId: updateId}
}

// TestTelegramRouteResponse needs tdlib installed to run. see: https://tdlib.github.io/td/build.html
func TestTelegramRouteResponse(t *testing.T) {
	var dates = lo.Map(timetable_gen.Timetable.Seasons,
		func(item timetable.Season, index int) time.Time {
			return item.Start
		})

	for _, _date := range dates {
		t.Run(_date.Format("2006-01-02"), func(t *testing.T) {
			testTelegramRouteOnDate(t, _date)
		})
	}
}

func testTelegramRouteOnDate(t *testing.T, _date time.Time) {
	_app, err := app.NewApp(date.FixedDate(_date))
	assert.NoError(t, err)
	mockTgClient := &mock_server.MockCustomTgClient{}
	serverStopped := make(chan struct{})
	defer func() { <-serverStopped }()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start server
	go func() {
		err := server.RunServer(ctx, _config, _app, server.WithCustomTgClient(mockTgClient))
		assert.NotErrorIs(t, err, http.ErrServerClosed)
		close(serverStopped)
	}()
	assert.Eventually(t, func() bool {
		conn, err := net.Dial("tcp", Url)
		if conn != nil {
			_ = conn.Close()
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
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, response.StatusCode, http.StatusOK)
	})
	// unknown station request — must run before "start message" to avoid stale sendMessage expectations
	// consuming the Run callback before it gets a chance to match
	t.Run("unknown station request", func(t *testing.T) {
		request := nextUpdateRequest()
		request.Message = &gotgbot.Message{
			Text: "Berlin, London",
			From: &gotgbot.User{
				LanguageCode: "en",
			},
		}

		requestRaw := lo.Must(json.Marshal(request))
		var capturedParams map[string]string
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "sendMessage", mock.Anything, mock.Anything, mock.Anything).
			Run(func(_ context.Context, _, _ string, params map[string]string, _ map[string]gotgbot.FileReader, _ *gotgbot.RequestOpts) {
				capturedParams = params
			}).
			Return([]byte("{}"), nil)
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
		assert.NoError(t, err)
		require.NotNil(t, response)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, model_render.StationDoesNotExistMessageMap[language.English], capturedParams["text"])
		assert.Contains(t, capturedParams["reply_markup"], model_render.GoogleMapWithAllStations)
	})
	// start message
	t.Run("start message", func(t *testing.T) {
		for _, languageTag := range model_render.SupportedLanguages {
			_language := languageTag.String()
			t.Run(_language, func(t *testing.T) {
				request := nextUpdateRequest()
				request.Message = &gotgbot.Message{
					Text: "/start",
					From: &gotgbot.User{
						LanguageCode: _language,
					},
				}

				requestRaw := lo.Must(json.Marshal(request))
				mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "sendMessage", mock.Anything, mock.Anything, mock.Anything).
					Return([]byte("{}"), nil)
				response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
				assert.NoError(t, err)
				require.NotNil(t, response)
				defer func() { _ = response.Body.Close() }()
				assert.Equal(t, response.StatusCode, http.StatusOK)
				t.Log("telegram response: ", lo.LastOrEmpty(mockTgClient.Calls).Arguments)
			})
		}
	})
	// timetable request
	t.Run("timetable request", func(t *testing.T) {
		request := nextUpdateRequest()
		request.Message = &gotgbot.Message{
			Text: "nikshich, bar",
			From: &gotgbot.User{
				LanguageCode: "en",
			},
		}

		requestRaw := lo.Must(json.Marshal(request))
		//mockTgClient.EXPECT().TimeoutContext(mock.Anything).Return(context.WithCancel(ctx))
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, TelegramApiToken, "sendMessage", mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("{}"), nil)
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
		assert.NoError(t, err)
		require.NotNil(t, response)
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, response.StatusCode, http.StatusOK)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
	// callback update
	t.Run("update callback", func(t *testing.T) {
		request := nextUpdateRequest()
		request.CallbackQuery = &gotgbot.CallbackQuery{
			Id:      "someid",
			Data:    "0|Niksic|Bar|2000-01-01",
			Message: &gotgbot.Message{},
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
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
	// callback update on the same date
	t.Run("update callback on the same date", func(t *testing.T) {
		request := nextUpdateRequest()
		request.CallbackQuery = &gotgbot.CallbackQuery{
			Id:      "someid",
			Data:    "0|Niksic|Bar|" + _date.Format("2006-01-02"),
			Message: &gotgbot.Message{},
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
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
	// callback reverse
	t.Run("reverse route callback", func(t *testing.T) {
		request := nextUpdateRequest()
		request.CallbackQuery = &gotgbot.CallbackQuery{
			Id:      "someid",
			Data:    "1|Niksic|Bar",
			Message: &gotgbot.Message{},
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
		defer func() { _ = response.Body.Close() }()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		t.Log("telegram response: ", mockTgClient.Calls[1].Arguments)
	})
}
