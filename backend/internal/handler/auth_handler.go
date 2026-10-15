package handler

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Torrahp/todo-backend/internal/repository"
	"github.com/Torrahp/todo-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

type credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req credentials
	if err := c.ShouldBindJSON(&req); err != nil { //จับคู่เข้าตัวJSON และ return ออกมาหลาย HTTP code
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.Service.Register(req.Email, req.Password)
	if errors.Is(err, repository.ErrEmailTaken) {
		c.JSON(409, gin.H{"error": "email already registered"})
		return
	}
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		c.JSON(400, gin.H{"error": "Password is too long"})
		return
	}
	if errors.Is(err, service.ErrInvalidCredentials) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to register"})
		return
	}
	c.JSON(201, gin.H{"message": "Registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req credentials
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Service.Login(req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to login"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
