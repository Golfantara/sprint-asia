package dtos

import "time"

type InputTask struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Deadline    string `json:"deadline,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type InputSubtask struct {
	ID          int       `json:"id"`
	TaskID      int       `json:"task_id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}