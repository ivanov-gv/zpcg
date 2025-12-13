package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

type BusinessConnection struct {
	Filter   filters.BusinessConnection
	Response Response
}

func (bc BusinessConnection) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	return ctx.BusinessConnection != nil && bc.Filter(ctx.BusinessConnection)
}

func (bc BusinessConnection) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return bc.Response(b, ctx)
}

func (bc BusinessConnection) Name() string {
	return fmt.Sprintf("businessconnection_%p", bc.Response)
}
