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

	mock_server "github.com/ivanov-gv/zpcg/gen/mocks/server"
	"github.com/ivanov-gv/zpcg/internal/app"
	"github.com/ivanov-gv/zpcg/internal/config"
	"github.com/ivanov-gv/zpcg/internal/server"
)

const (
	TelegramApiToken  = "TelegramApiToken"
	Port              = "8080" // TODO: take random open
	Url               = "localhost:" + Port
	HttpServerAddress = "http://" + Url
	Environment       = "TestEnvironment"
)

var _config = config.Config{
	TelegramApiToken: TelegramApiToken,
	Port:             Port,
	Environment:      Environment,
}

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
		defer conn.Close()
		return err == nil
	}, time.Second, 100*time.Millisecond)

	// run tests
	// send empty message update
	t.Run("message update", func(t *testing.T) {
		const request = "{}"
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer([]byte(request)))
		assert.NoError(t, err)
		assert.Equal(t, response.StatusCode, http.StatusOK)
	})
	t.Run("message update", func(t *testing.T) {
		request := gotgbot.Update{Message: &gotgbot.Message{Text: "nikshich, bar"}}
		requestRaw := lo.Must(json.Marshal(request))
		mockTgClient.EXPECT().TimeoutContext(mock.Anything).Return(context.WithCancel(ctx))
		mockTgClient.EXPECT().RequestWithContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("{}"), nil)
		response, err := http.Post(HttpServerAddress, "application/json", bytes.NewBuffer(requestRaw))
		assert.NoError(t, err)
		assert.Equal(t, response.StatusCode, http.StatusOK)
	})
}
