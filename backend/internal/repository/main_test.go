//go:build integration

package repository

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Torrahp/todo-backend/internal/config"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	if !strings.HasSuffix(os.Getenv("DB_NAME"), "_test") {
		log.Fatal("DB_NAME must end with _test (refusing to wipe a non-test database)")
	}
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("connect test db: %v", err)
	}
	testDB = db

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func resetDB(t *testing.T) {
	t.Helper()
	if _, err := testDB.Exec(`TRUNCATE users, todos RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("reset db: %v", err)
	}
}

func createUser(t *testing.T, email string) int {
	t.Helper()
	repo := &UserRepository{DB: testDB}
	if err := repo.CreateUser(email, "not-a-real-hash"); err != nil {
		t.Fatalf("create user: %v", err)
	}
	u, err := repo.GetByEmail(email)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	return u.ID
}
