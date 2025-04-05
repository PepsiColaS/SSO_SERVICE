package postgres

import (
	"context"
	"fmt"
	"sso/internal/config"

	"github.com/jackc/pgx/v4"
)

type PostgresClient struct {
	conn *pgx.Conn
}

func ConnectToDataBase(cfg *config.Config) *PostgresClient {
	dataBaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.GRPC.User, cfg.GRPC.Password, cfg.GRPC.Host, cfg.GRPC.Db)
	conn, err := pgx.Connect(context.Background(), dataBaseURL)
	if err != nil {
		panic(err)
	}
	return &PostgresClient{conn}
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

func (pc *PostgresClient) CreateDataBase() {
	query := `
			CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			pass_hash TEXT  NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_login ON users (login);

		CREATE TABLE IF NOT EXISTS apps (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			secret TEXT NOT NULL UNIQUE
		); `
	_, err := pc.conn.Exec(context.Background(), query)
	if err != nil {
		panic(err)
	}
}
