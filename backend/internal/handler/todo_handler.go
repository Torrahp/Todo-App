package handler

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Torrahp/todo-backend/internal/repository"
)

type TodoHandler struct {
	Repo *repository.TodoRepository
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := c.GetInt("userID")
	todos, err := h.Repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todos"})
		return
	}
	c.JSON(200, todos)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required,max=255"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		c.JSON(400, gin.H{"error": "title cannot be empty"})
		return
	}
	userID := c.GetInt("userID")
	todo, err := h.Repo.Create(userID, title)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create todo"})
		return
	}
	c.JSON(201, todo) //created
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetInt("userID")

	todo, err := h.Repo.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, gin.H{"error": "todo not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todo"})
		return
	}
	if todo.UserID != userID {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	var req struct {
		Title  *string `json:"title" binding:"omitempty,max=255"`
		IsDone *bool   `json:"is_done"`
		//ที่ต้องเป็น pointer เพราะ Go ถ้าไม่ส่งค่ามาจะเป็นfalse ถ้าเป็นpointer จะเป็นnil ทำให้เรารู้ว่า user ไม่ได้ส่งค่าอะไรมา
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			c.JSON(400, gin.H{"error": "title cannot be empty"})
			return
		}
		todo.Title = title
	}
	if req.IsDone != nil {
		todo.IsDone = *req.IsDone
		// ถ้าเป้น ค่า true, false อันนี้ปกติ เพราะ userส่งมาเอง
		// ถ้าเป็น nil  เพราะ user ไม่ได้ส่งค่าอะไรมา เราจะไม่เปลี่ยนค่าเดิม
	}

	updated, err := h.Repo.Update(todo.ID, todo.Title, todo.IsDone)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update todo"})
		return
	}
	c.JSON(200, updated)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetInt("userID")

	todo, err := h.Repo.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, gin.H{"error": "todo not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todo"})
		return
	}
	if todo.UserID != userID {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	if err := h.Repo.Delete(todo.ID); err != nil {
		c.JSON(500, gin.H{"error": "failed to delete todo"})
		return
	}
	c.Status(204) //delete สำเร็จจะส่งแบบนี่เพราะ ไม่ต้องส่งค่าอะไรกลับ
}
