package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Torrahp/todo-backend/internal/config"
	"github.com/Torrahp/todo-backend/internal/handler"
	"github.com/Torrahp/todo-backend/internal/middleware"
	"github.com/Torrahp/todo-backend/internal/repository"
	"github.com/Torrahp/todo-backend/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no ,env file found, using system env")
	}
	if config.GetEnv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("database connect failed: %v", err) //error แล้วจบexits(1)
	}
	defer db.Close() // ปิด connnection หลังจาก main func จบ

	todoRepo := &repository.TodoRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}
	authService := &service.AuthService{Repo: userRepo}
	todoHandler := &handler.TodoHandler{Repo: todoRepo}
	authHandler := &handler.AuthHandler{Service: authService}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.GetEnv("CORS_ORIGIN", "http://localhost:5173")},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}))
	// ทำให้หน้าเว็บที่อยู่ใน AllowOrigins สามารถเรียก API ของเราได้ และสามารถส่ง Header Authorization และ Content-Type ได้
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Status": "Ok"})
	})

	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	todos := r.Group("/todos")
	todos.Use(middleware.AuthRequired())
	todos.GET("", todoHandler.GetTodos)
	todos.POST("", todoHandler.CreateTodo)
	todos.PATCH("/:id", todoHandler.UpdateTodo)
	todos.DELETE("/:id", todoHandler.DeleteTodo)

	if err := r.Run(":" + config.GetEnv("PORT", "8080")); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
