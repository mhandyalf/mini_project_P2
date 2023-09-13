package main

import (
	"mini_project_p2/router"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	router.SetupRouter(e)

	e.Logger.Fatal(e.Start(":8080"))

}
