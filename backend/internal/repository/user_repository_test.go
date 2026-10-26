//go:build integration

package repository

import (
	"errors"
	"testing"
)

func TestCreateUser_DuplicateEmail(t *testing.T) {
	resetDB(t)
	repo := &UserRepository{DB: testDB}

	if err := repo.CreateUser("a@example.com", "hash"); err != nil {
		t.Fatalf("first create: %v", err)
	}
	err := repo.CreateUser("a@example.com", "hash")
	if !errors.Is(err, ErrEmailTaken) {
		t.Errorf("err = %v, want ErrEmailTaken", err)
	}
}
