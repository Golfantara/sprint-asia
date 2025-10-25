package handler

import (
	"case-study-2/helpers"
	helper "case-study-2/helpers"
	"fmt"
	"strconv"

	"case-study-2/features/task"
	"case-study-2/features/task/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service task.Usecase
}

func New(service task.Usecase) task.Handler {
	return &controller {
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetTasks() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)
		
		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		tasks := ctl.service.FindAll(page, size)

		if tasks == nil {
			return ctx.JSON(404, helper.Response("There is No Tasks!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": tasks,
		}))
	}
}


func (ctl *controller) TaskDetails() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		taskID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		task := ctl.service.FindByID(taskID)

		if task == nil {
			return ctx.JSON(404, helper.Response("Task Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": task,
		}))
	}
}

func (ctl *controller) CreateTask() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputTask{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		task := ctl.service.Create(input)

		if task == nil {
			return ctx.JSON(500, helper.Response("Something went Wrong!", nil))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": task,
		}))
	}
}

func (ctl *controller) UpdateTask() echo.HandlerFunc {
	return func (ctx echo.Context) error {
		input := dtos.InputTask{}

		taskID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		task := ctl.service.FindByID(taskID)

		if task == nil {
			return ctx.JSON(404, helper.Response("Task Not Found!"))
		}
		
		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		update := ctl.service.Modify(input, taskID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Task Success Updated!"))
	}
}

func (ctl *controller) DeleteTask() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		taskID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		task := ctl.service.FindByID(taskID)

		if task == nil {
			return ctx.JSON(404, helper.Response("Task Not Found!"))
		}

		delete := ctl.service.Remove(taskID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Task Success Deleted!", nil))
	}
}

func (ctl *controller) GetByTaskID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		taskID, err := strconv.Atoi(ctx.Param("task_id"))
		fmt.Println("cek ",taskID)
		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}
		data := ctl.service.FindByTaskID(taskID)
		return ctx.JSON(200, helpers.Response("Success", map[string]any{"data": data}))
	}
}

func (ctl *controller) CreateSubtask() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var input dtos.InputSubtask
		ctx.Bind(&input)

		data := ctl.service.CreateSubtask(input)
		if data == nil {
			return ctx.JSON(500, helpers.Response("Failed to create subtask"))
		}
		return ctx.JSON(200, helpers.Response("Success", map[string]any{"data": data}))
	}
}

func (ctl *controller) GetHistoryTasks() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide valid pagination"))
		}

		tasks := ctl.service.FindHistory(page, size)
		if len(tasks) == 0 {
			return ctx.JSON(404, helper.Response("No completed tasks found"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": tasks,
		}))
	}
}
func (ctl *controller) UpdateSubtask() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		subtaskID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		var input dtos.InputSubtask
		ctx.Bind(&input)

		success := ctl.service.ModifySubtask(input, subtaskID)
		if !success {
			return ctx.JSON(500, helpers.Response("Failed to update subtask"))
		}
		return ctx.JSON(200, helpers.Response("Subtask successfully updated!"))
	}
}
