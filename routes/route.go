package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"Users/controller"
	"Users/db"
	mid "Users/middleware"
	"Users/repository"
	"Users/service"

)

func Init() error {
	e := echo.New()

	db, err := db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	controller := controller.NewController(service)

	task := e.Group("/task")
	task.Use(mid.ValidateToken)

	// Routes
	e.GET("", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "Application is Running",
		})
	})

	task.POST("", controller.CreateTasksController)
	task.DELETE("/:id", controller.DeleteTasksController)
	task.GET("", controller.GetAlltaskController)
	task.GET("/:id", controller.GetTaskById)
	task.PUT("/:id", controller.UpdateTaskController)
	task.DELETE("", controller.BulkDeleteTask)
	task.GET("/count", controller.CountTask)
	task.POST("/search", controller.SearchTasksFormController)

	e.POST("/login", controller.Login)
	e.POST("/forget-password", controller.ForgotPasswordHandler)
	e.POST("/reset-password", controller.ResetPassword)
	e.POST("/logout", controller.Logout)
	e.POST("/register", controller.RegisterController)

	return e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
