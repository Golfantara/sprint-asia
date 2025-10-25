package task

import "time"

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"column:email;type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"column:password;type:text;not null" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Task struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      int        `gorm:"column:user_id" json:"user_id"`
	Title       string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	IsCompleted bool       `gorm:"column:is_completed;default:false" json:"is_completed"`
	Deadline    *time.Time `gorm:"column:deadline" json:"deadline,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Subtasks []SubTask `gorm:"foreignKey:TaskID;references:ID" json:"subtasks"`
}


func (Task) TableName() string {
	return "tasks"
}

type SubTask struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID      int       `gorm:"column:task_id" json:"task_id"`
	Title       string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	IsCompleted bool      `gorm:"column:is_completed;default:false" json:"is_completed"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SubTask) TableName() string {
	return "subtasks"
}
