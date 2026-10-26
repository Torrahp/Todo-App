package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Torrahp/todo-backend/internal/model"
)

type fakeUserStore struct {
	users map[string]model.User
}

func (f *fakeUserStore) CreateUser(email string, passwordHash string) error {
	f.users[email] = model.User{ID: len(f.users) + 1, Email: email, PasswordHash: passwordHash}
	return nil
}

func (f *fakeUserStore) GetByEmail(email string) (model.User, error) {
	u, ok := f.users[email]
	if !ok {
		return model.User{}, sql.ErrNoRows
	}
	return u, nil
}

func TestRegister(t *testing.T) {
	store := &fakeUserStore{users: map[string]model.User{}}
	svc := &AuthService{Repo: store}

	if err := svc.Register("a@example.com", "password123"); err != nil {
		t.Fatalf("register: %v", err)
	}
	hash := store.users["a@example.com"].PasswordHash
	if hash == "password123" {
		t.Fatal("password stroed in plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("password123")); err != nil {
		t.Errorf("stored hash does not match password: %v", err)
	}
}

func TestLogin(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	svc := &AuthService{Repo: &fakeUserStore{users: map[string]model.User{}}}
	if err := svc.Register("a@example.com", "password123"); err != nil {
		t.Fatalf("register: %v", err)
	}

	t.Run("wrong password", func(t *testing.T) {
		_, err := svc.Login("a@example.com", "wrong-password")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("err = %v, want ErrInvaildCredentials", err)
		}
	})

	t.Run("unknown email", func(t *testing.T) {
		_, err := svc.Login("nobodynobodybonju@example.com", "password123")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("err = %v, want ErrInvaildCredentials", err)
		}
	})

	t.Run("correct password return token with user_id", func(t *testing.T) {
		tokenStr, err := svc.Login("a@example.com", "password123")
		if err != nil {
			t.Fatalf("login: %v", err)
		}
		token, err := jwt.Parse(tokenStr, func(*jwt.Token) (any, error) {
			return []byte("test-secret"), nil
		})
		if err != nil || !token.Valid {
			t.Fatalf("token invaild: %v", err)
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["user_id"] != float64(1) {
			t.Errorf("user_id = %v, want 1", claims["user_id"])
		}
	})
}
