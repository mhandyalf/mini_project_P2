package router

import (
	"mini_project_p2/database"
	"mini_project_p2/handlers"
	"mini_project_p2/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo) {

	auth := handlers.NewAuth(database.InitDB())

	e.POST("/register", auth.Register)
	e.POST("/login", auth.Login)

	e.GET("/books", auth.GetAllBooks, middleware.JWTAuth)
	e.POST("/rent", auth.RentBook, middleware.JWTAuth)
	e.PUT("/update", auth.UpdateBook, middleware.JWTAuth)
	e.DELETE("/delete", auth.DeleteBook, middleware.JWTAuth)

}
