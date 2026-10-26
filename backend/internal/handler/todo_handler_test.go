package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Torrahp/todo-backend/internal/model"
	"github.com/gin-gonic/gin"
)

type fakeTodoStore struct {
	todos   map[int]model.Todo
	failGet bool
}

func newFakeStore() *fakeTodoStore {
	return &fakeTodoStore{todos: map[int]model.Todo{
		1: {ID: 1, UserID: 1, Title: "mine"},
		2: {ID: 2, UserID: 2, Title: "not mine"},
	}}
}

func (f *fakeTodoStore) GetAllByUser(userID int) ([]model.Todo, error) {
	out := []model.Todo{}
	for _, t := range f.todos {
		if t.UserID == userID {
			out = append(out, t)
		}
	}
	return out, nil
}

func (f *fakeTodoStore) GetByID(id int) (model.Todo, error) {
	if f.failGet {
		return model.Todo{}, errors.New("db down")
	}
	t, ok := f.todos[id] // เป็นการเช็คว่าค่า f มีอยู่ในtodosไหม
	if !ok {
		return model.Todo{}, sql.ErrNoRows
	}
	return t, nil
}

func (f *fakeTodoStore) Create(userID int, title string) (model.Todo, error) {
	t := model.Todo{
		ID: 100, UserID: userID, Title: title,
	}
	f.todos[t.ID] = t
	return t, nil
}

func (f *fakeTodoStore) Update(id int, title string, isDone bool) (model.Todo, error) {
	t, ok := f.todos[id]
	if !ok {
		return model.Todo{}, sql.ErrNoRows
	}
	t.Title = title
	t.IsDone = isDone
	f.todos[id] = t
	return t, nil
}

func (f *fakeTodoStore) Delete(id int) error {
	_, ok := f.todos[id]
	if !ok {
		return sql.ErrNoRows
	}
	delete(f.todos, id)

	return nil
}

func setupRouter(store TodoStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := &TodoHandler{Repo: store}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	r.GET("/todos", h.GetTodos)
	r.POST("/todos", h.CreateTodo)
	r.PATCH("/todos/:id", h.UpdateTodo)
	r.DELETE("/todos/:id", h.DeleteTodo)
	return r
}

func TestTodoHandler_StatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		failGet    bool
		wantStatus int
	}{
		{"list", "GET", "/todos", "", false, 200},
		{"create ok", "POST", "/todos", `{"title":"buy milk"}`, false, 201},
		{"create without title", "POST", "/todos", `{}`, false, 400},
		{"create with blank title", "POST", "/todos", `{"title":"   "}`, false, 400},
		{"update invalid id", "PATCH", "/todos/abc", `{"title":"x"}`, false, 400},
		{"update not found", "PATCH", "/todos/99", `{"title":"x"}`, false, 404},
		{"update other user's todo", "PATCH", "/todos/2", `{"title":"x"}`, false, 403},
		{"update empty title", "PATCH", "/todos/1", `{"title":""}`, false, 400},
		{"update ok", "PATCH", "/todos/1", `{"is_done":true}`, false, 200},
		{"delete invalid id", "DELETE", "/todos/abc", "", false, 400},
		{"delete not found", "DELETE", "/todos/99", "", false, 404},
		{"delete other user's todo", "DELETE", "/todos/2", "", false, 403},
		{"delete when db fails", "DELETE", "/todos/1", "", true, 500},
		{"delete ok", "DELETE", "/todos/1", "", false, 204},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := newFakeStore()
			store.failGet = tc.failGet
			r := setupRouter(store)

			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", w.Code, tc.wantStatus, w.Body.String())
			}
		})
	}
}

func TestCreateTodo_SavesTrimmedTitle(t *testing.T) {
	store := newFakeStore()
	r := setupRouter(store)

	req := httptest.NewRequest("POST", "/todos", strings.NewReader(`{"title":"  buy milk  "}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)

	if got := store.todos[100].Title; got != "buy milk" {
		t.Errorf("saved title = %q, want %q", got, "buy milk")
	}
}

func TestGetTodos_ReturnsOnlyOwnTodos(t *testing.T) {
	r := setupRouter(newFakeStore())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/todos", nil))

	var got []model.Todo
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 1 || got[0].UserID != 1 {
		t.Errorf("got %+v, want only todos of user 1", got)
	}
}
