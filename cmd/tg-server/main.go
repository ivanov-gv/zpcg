package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"zpcg/internal/app"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI("6813690978:AAHACdR-uUs3vjpu07EwX-ppRvUw3YUDwaI")
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// create app
	_app, err := app.NewApp("./resources/timetable.gob")
	if err != nil {
		log.Panic(err)
	}

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, _app, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, _app *app.App, updates tgbotapi.UpdatesChannel) {
	for {
		if ctx.Err() != nil {
			return
		}
		update := <-updates
		if update.Message == nil {
			continue
		}
		handleMessage(_app, update.Message)
	}
}

func handleMessage(_app *app.App, message *tgbotapi.Message) {
	if message.From == nil {
		return
	}
	log.Printf("%s %s wrote %s", message.From.FirstName, message.From.UserName, message.Text)

	// generate message
	var msgText, parseMode string
	// got command - send start message
	if strings.HasPrefix(message.Text, "/") {
		msgText, parseMode = _app.StartMessage()
		// got two stations - send timetable
	} else if originStation, destinationStation, found := strings.Cut(message.Text, ","); found {
		msgText, parseMode = _app.GenerateRoute(originStation, destinationStation)
		// error otherwise
	} else {
		msgText, parseMode = _app.ErrorMessage()
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ParseMode = parseMode
	_, err := bot.Send(msg)
	log.Println("got error: ", err)
}
