package main

import (
	"context"
	"fin_notifications/cmd/commands"
	"fin_notifications/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const defaultEnvFilePath = ".env"

func main() {
	cfg, err := config.Parse(defaultEnvFilePath)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel()
	}()

	err = commands.ReadFromQueue(ctx, cfg)
	if err != nil {
		slog.Error("Error reading from queue", "error", err)
		return
	}

}
