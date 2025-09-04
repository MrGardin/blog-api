package repository

import (
	"database/sql"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CheckConnection() error {
	return r.db.Ping()
}

func (r *UserRepository) GetWelcomeMessage() (string, error) {
	var msg string
	err := r.db.QueryRow("SELECT 'База данных работает'::text").Scan(&msg)
	if err != nil {
		log.Fatalf("Не удалось отправить запрос: %v", err)
	}
	return msg, err
}
