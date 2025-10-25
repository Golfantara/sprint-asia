package dtos

import "time"
type ResTask struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	IsCompleted bool          `json:"is_completed"`
	Deadline    *time.Time    `json:"deadline,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Progress    int           `json:"progress"`
	IsDue       bool          `json:"is_due"`
	Subtask     []ResSubtask  `json:"subtasks"`
}

type ResSubtask struct {
	ID          int       `json:"id"`
	TaskID      int       `json:"task_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
