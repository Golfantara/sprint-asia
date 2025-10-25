package task

import (
	"case-study-2/features/task/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Task
	Insert(newTask Task) int64
	SelectByID(taskID int) *Task
	Update(task Task) int64
	DeleteByID(taskID int) int64
	GetByTaskID(taskID int) []SubTask
	InsertSubtask(newSubtask SubTask) int64
	GetHistoryTasks(page, size int) []Task
	UpdateSubtask(subtask SubTask) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResTask
	FindByID(taskID int) *dtos.ResTask
	Create(newTask dtos.InputTask) *dtos.ResTask
	Modify(taskData dtos.InputTask, taskID int) bool
	Remove(taskID int) bool
	FindByTaskID(taskID int) []dtos.ResSubtask
	CreateSubtask(input dtos.InputSubtask) *dtos.ResSubtask
	ModifySubtask(input dtos.InputSubtask, subtaskID int) bool
	FindHistory(page, size int) []dtos.ResTask
}

type Handler interface {
	GetTasks() echo.HandlerFunc
	TaskDetails() echo.HandlerFunc
	CreateTask() echo.HandlerFunc
	UpdateTask() echo.HandlerFunc
	DeleteTask() echo.HandlerFunc
	GetByTaskID() echo.HandlerFunc
	CreateSubtask() echo.HandlerFunc
	UpdateSubtask() echo.HandlerFunc 
	GetHistoryTasks() echo.HandlerFunc
}

