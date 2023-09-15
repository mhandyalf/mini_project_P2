package main

import (
	"mini_project_p2/router"

	_ "mini_project_p2/docs"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

// @title Handy Library API
// @version 1.0
// @description This is Handy Library API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v2

func main() {

	e := echo.New()

	router.SetupRouter(e)

	e.Logger.Fatal(e.Start(":8080"))

}
