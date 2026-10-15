package repository

import (
	"database/sql"
	"errors" // หากอยากใช้ error ของเรา

	"github.com/Torrahp/todo-backend/internal/model"
	"github.com/lib/pq"
)

var ErrEmailTaken = errors.New("email already registered")

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(email, passwordHash string) error {
	_, err := r.DB.Exec(`INSERT INTO users (email, password_hash) VALUES ($1, $2)`,
		email, passwordHash)
	var pqErr *pq.Error                                  // เตรียชนิดตัวแปร
	if errors.As(err, &pqErr) && pqErr.Code == "23505" { //errors.AS เป็นการแปลง จาก interface error เป็น
		return ErrEmailTaken
	}
	return err
}

func (r *UserRepository) GetByEmail(email string) (model.User, error) {
	var u model.User
	err := r.DB.QueryRow(`SELECT id, email, password_hash, created_at FROM users WHERE email=$1`,
		email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}
