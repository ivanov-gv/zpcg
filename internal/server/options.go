package server

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type options struct {
	botOpts *gotgbot.BotOpts
}

type ApplyOption func(opts options) options

func applyOptions(opts options, appliers ...ApplyOption) options {
	for _, applier := range appliers {
		opts = applier(opts)
	}
	return opts
}

type CustomTgClient interface {
	gotgbot.BotClient
}

func WithCustomTgClient(httpClient CustomTgClient) ApplyOption {
	return func(opts options) options {
		opts.botOpts.BotClient = httpClient
		return opts
	}
}
