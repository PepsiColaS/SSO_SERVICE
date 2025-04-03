package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/database/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, port int, storage *postgres.PostgresClient) *App {
	grpcApp := grpcapp.New(log, port, storage)

	// authService := auth.New(log, storage, storage, storage, tokenTTL)
	// grpcApp := grpcapp.New(log, port, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
