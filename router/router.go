package router

import (
	"mini_project_p2/database"
	"mini_project_p2/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo) {
	auth := handlers.NewAuth(database.InitDB())

	e.POST("/register", auth.Register)
	e.POST("/login", auth.Login)

}
