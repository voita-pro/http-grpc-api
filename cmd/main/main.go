package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/voita-pro/http-grpc-api/internal/app"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("can't start application")
	}
	err = a.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("can't run application")
	}
}
