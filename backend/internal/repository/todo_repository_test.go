//go:build integration

package repository

import (
	"database/sql"
	"errors"
	"testing"
)

func TestTodoRepository_CRUD(t *testing.T) {
	resetDB(t)
	repo := &TodoRepository{DB: testDB}
	userID := createUser(t, "a@example.com")

	created, err := repo.Create(userID, "buy milk")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == 0 || created.Title != "buy milk" || created.IsDone {
		t.Fatalf("unexpected todo: %+v", created)
	}

	updated, err := repo.Update(created.ID, "buy oat milk", true)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "buy oat milk" || !updated.IsDone {
		t.Errorf("update not saved: %+v", updated)
	}

	if err := repo.Delete(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.GetByID(created.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("after delete err = %v, want sql.ErrNoRows", err)
	}
}

func TestGetAllByUser_OnlyOwnTodosNewestFirst(t *testing.T) {
	resetDB(t)
	repo := &TodoRepository{DB: testDB}
	alice := createUser(t, "alice@example.com")
	bob := createUser(t, "bob@example.com")

	repo.Create(alice, "first")
	repo.Create(alice, "second")
	repo.Create(bob, "bob's todo")

	todos, err := repo.GetAllByUser(alice)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("got %d todos, want 2", len(todos))
	}
	if todos[0].Title != "second" {
		t.Errorf("first item = %q, want newest (\"second\")", todos[0].Title)
	}
}

func TestDeleteUser_CascadesTodos(t *testing.T) {
	resetDB(t)
	repo := &TodoRepository{DB: testDB}
	userID := createUser(t, "a@example.com")
	todo, _ := repo.Create(userID, "will be gone")

	if _, err := testDB.Exec(`DELETE FROM users WHERE id=$1`, userID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if _, err := repo.GetByID(todo.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("todo still exists after user deleted (err = %v)", err)
	}
}
