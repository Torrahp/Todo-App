package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // เวลาดึงข้อมูล จะไม่ได้ filed นี่ออกมา
	CreatedAt    time.Time `json:"created_at"`
}
