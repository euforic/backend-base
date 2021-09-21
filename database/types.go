package database

import (
	"time"
)

// Todo contains all the information associated to a user.
type Todo struct {
	Id        string     `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Author    string     `json:"author"`
	IsDone    *bool      `json:"is_done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// NewTodo returns an Todo instance.
func NewTodo() *Todo {
	return &Todo{}
}
