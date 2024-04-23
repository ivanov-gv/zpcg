package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ivanov-gv/zpcg/internal/app"
	"github.com/ivanov-gv/zpcg/internal/config"
	"github.com/ivanov-gv/zpcg/internal/server"
)

func main() {
	const logfmt = "main: "
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// config
	_config, err := config.Load()
	if err != nil {
		log.Fatal().Err(fmt.Errorf(logfmt+"can't load config: %w", err)).Send()
	}
	// app
	_app, err := app.NewApp()
	if err != nil {
		log.Fatal().Err(fmt.Errorf(logfmt+"app.NewApp: %w", err)).Send()
	}
	// server
	err = server.RunServer(ctx, _config, _app)
	if err != nil {
		log.Fatal().Err(fmt.Errorf(logfmt+"server.RunServer: %w", err)).Send()
	}
}
