package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	os.Exit(run()) //nolint:forbidigo // main entry point requires os.Exit
}

func run() int {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelCtx()

	cfg, err := NewConfigFromEnv()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		return 1
	}
	emailSender := NewEmailSender(cfg)
	defer emailSender.Shutdown()
	h := NewHandler(cfg, emailSender)
	srv := NewServer(cfg, h)
	if err := srv.Start(); err != nil {
		slog.Error("server failed to start", "err", err)
		return 1
	}
	defer srv.Stop()

	<-ctx.Done()
	return 0
}
