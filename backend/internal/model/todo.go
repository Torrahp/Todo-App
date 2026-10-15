package model

import "time"

type Todo struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}

//filed_Name(use in go) TYPE คือ struct tag ไว้Encode เป็นJson
