package utils

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"websocket-chat/internal/store"

	"github.com/joho/godotenv"
)

func InitialiseDb() (*store.SQLStore, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	url := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	fullURL := fmt.Sprintf("%s?authToken=%s", url, authToken)
	db, err := sql.Open("libsql", fullURL)

	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", url, err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return store.NewSQLStore(db), nil
}
