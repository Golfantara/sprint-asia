package main

import (
	"case-study-2/features/task"
	"case-study-2/routes"
	"case-study-2/utils"
	"fmt"
	"net/http"

	ch "case-study-2/features/task/handler"
	cr "case-study-2/features/task/repository"
	cu "case-study-2/features/task/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.Tasks(e, TaskHandler())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello anjay mabar!")
	})
	e.Start(fmt.Sprintf(":8000"))
}

func TaskHandler() task.Handler {

	db := utils.InitDB()

	repo := cr.New(db)
	cc :=	cu.New(repo)
	return ch.New(cc)
}