package service

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Torrahp/todo-backend/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	Repo *repository.UserRepository
}

func (s *AuthService) Register(email, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// จะสุ่มค่าsalt ออกมาทำให้ รหัสผ่านที่เหมือนกัน มีpassword hash ไม่เหมือนกัน
	//จากนั่นเอาไป encode โดนมีค่าcost เพื่อทำให้ เข้ารหัสได้ช้า ป้องกันการโจมตี
	if err != nil {
		return err
	}
	return s.Repo.CreateUser(strings.ToLower(email), string(hash))
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetByEmail(strings.ToLower(email))
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	return token.SignedString([]byte(jwtSecret))
}
