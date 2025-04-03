package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/database/postgres"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	storage, err := postgres.ConnectToDataBase(cfg, log)
	if err != nil {
		panic(err)
	}

	application := app.New(log, cfg.GRPC.ApiPort, storage)
	go application.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	storage.Close()
	log.Info("storage stopped")
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
