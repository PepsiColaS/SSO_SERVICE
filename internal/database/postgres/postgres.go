package postgres

import (
	"context"
	"fmt"
	"sso/internal/config"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
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

func (pc *PostgresClient) Login(login string, password string) (uID int64, err error) {
	passHash, err := pc.User(login)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(passHash, []byte(password)); err != nil {
		return 0, err
	}
	var id int64
	err = pc.conn.QueryRow(context.Background(), "SELECT id FROM users WHERE login = $1", login).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("User not found")
		}
		fmt.Println(err, '3')
		return 0, err
	}
	return id, nil
}

func (pc *PostgresClient) RegisterNewUser(login string, password string) (uID int64, err error) {
	var id int64
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	err = pc.conn.QueryRow(context.Background(), "INSERT INTO users (login, pass_hash) VALUES ($1, $2) RETURNING id", login, passHash).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (pc *PostgresClient) IsAdmin(uID int64) (isAdmin bool, err error) {
	return false, err
}

func (pc *PostgresClient) IsLibrarian(uID int64) (isLibrarian bool, err error) {
	return false, err
}

func (pc *PostgresClient) User(login string) (passHash []byte, err error) {
	var returnPassHash []byte
	err = pc.conn.QueryRow(context.Background(), "SELECT pass_hash FROM users WHERE login = $1", login).Scan(&returnPassHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}
	return returnPassHash, nil
}

func (pc *PostgresClient) CreateDataBase() {
	query := `
			CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			pass_hash TEXT  NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE
			isLibrarian BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE INDEX IF NOT EXISTS idx_login ON users (login);
		`
	_, err := pc.conn.Exec(context.Background(), query)
	if err != nil {
		panic(err)
	}
}
