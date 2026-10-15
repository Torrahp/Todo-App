package repository

import (
	"database/sql"

	"github.com/Torrahp/todo-backend/internal/model"
)

type TodoRepository struct {
	DB *sql.DB
}

func (r *TodoRepository) GetAllByUser(userID int) ([]model.Todo, error) {
	rows, err := r.DB.Query(
		`SELECT id, user_id, title, is_done, created_at FROM todos WHERE user_id=$1 ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []model.Todo{}
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepository) GetByID(id int) (model.Todo, error) {
	var t model.Todo
	err := r.DB.QueryRow(
		`SELECT id, user_id, title, is_done, created_at FROM todos WHERE id=$1`,
		id).Scan(&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt)
	return t, err
}

func (r *TodoRepository) Create(userID int, title string) (model.Todo, error) {
	var t model.Todo
	err := r.DB.QueryRow(
		`INSERT INTO todos (user_id, title) VALUES ($1, $2) RETURNING id, user_id, title, is_done, created_at`,
		userID, title).Scan(&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt)
	return t, err
}

func (r *TodoRepository) Update(id int, title string, isDone bool) (model.Todo, error) {
	var t model.Todo
	err := r.DB.QueryRow(
		`UPDATE todos SET title=$1, is_done=$2, updated_at=now() WHERE id=$3 RETURNING id, user_id, title, is_done, created_at`,
		title, isDone, id).Scan(&t.ID, &t.UserID, &t.Title, &t.IsDone, &t.CreatedAt)
	return t, err
}

func (r *TodoRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM todos WHERE id=$1`, id)
	return err
}
