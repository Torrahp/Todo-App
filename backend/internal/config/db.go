package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", GetEnv("DB_HOST"), GetEnv("DB_PORT"), GetEnv("DB_USER"), GetEnv("DB_PASSWORD"), GetEnv("DB_NAME"))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
