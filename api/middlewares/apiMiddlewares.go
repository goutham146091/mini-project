package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetApiMiddlewares(e *echo.Echo) {
	// this logs the webserver interaction
	e.Use((middleware.Logger()))
}
