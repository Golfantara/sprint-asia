package routes

import (
	"case-study-2/features/task"

	"github.com/labstack/echo/v4"
)

func Tasks(e *echo.Echo, handler task.Handler) {
	tasks := e.Group("/tasks")


	tasks.GET("", handler.GetTasks())
	tasks.GET("/history", handler.GetHistoryTasks())
	tasks.POST("", handler.CreateTask())
	tasks.GET("/:id", handler.TaskDetails())
	tasks.PUT("/:id", handler.UpdateTask())
	tasks.DELETE("/:id", handler.DeleteTask())

	subTasks := tasks.Group("/subtask")
	subTasks.GET("/:task_id", handler.GetByTaskID())
	subTasks.POST("", handler.CreateSubtask())
	subTasks.PUT("/:id", handler.UpdateSubtask())
}
