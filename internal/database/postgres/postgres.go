package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/config"

	"github.com/jackc/pgx/v4"
)

type PostgresClient struct {
	conn *pgx.Conn
}

func ConnectToDataBase(cfg *config.Config, log *slog.Logger) (*PostgresClient, error) {
	dataBaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.GRPC.User, cfg.GRPC.Password, cfg.GRPC.Host, cfg.GRPC.Db)
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), dataBaseURL)
	if err != nil {
		log.Error("Failed connet to database")
		return nil, err
	}
	log.Info("Successful connect with DB")
	return &PostgresClient{conn}, nil
}

func (pc *PostgresClient) Close() error {
	return pc.conn.Close(context.Background())
}

func (pc *PostgresClient) Login(login string, password string, appID int) (token string, err error) {
	return "fggf", err
}

func (pc *PostgresClient) RegisterNewUser(login string, password string) (userID int64, err error) {
	return 12, err
}

func (pc *PostgresClient) IsAdmin(uID int64) (isAdmin bool, err error) {
	return false, err
}
