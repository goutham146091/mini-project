package middlewares

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetAdminMiddlewares(adminGroup *echo.Group) {
	log.Println("jwt error")
	adminGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("gouthamospsm"),
	}))
}
