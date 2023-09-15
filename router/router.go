package router

import (
	"mini_project_p2/database"
	"mini_project_p2/handlers"
	"mini_project_p2/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo) {
	auth := handlers.NewAuth(database.InitDB())

	// Routes without JWT middleware
	e.POST("/register", auth.Register)
	e.POST("/login", auth.Login)

	// Routes with JWT middleware
	protected := e.Group("")
	protected.Use(middleware.JWTAuth)

	protected.GET("/books", auth.GetAllBooks)
	protected.POST("/rent", auth.RentBook)
	protected.PUT("/update", auth.UpdateBook)
	protected.DELETE("/delete", auth.DeleteBook)
	protected.POST("/topup", auth.TopUpDeposit)
	protected.POST("/return", auth.ReturnBook)
	protected.DELETE("/deletehistory", auth.DeleteRentalHistory)
}
