package api

import (
	"git.npcompete.com/OSPSM_Servers/src/api/handlers"

	"github.com/labstack/echo"
)

func CommonGroup(e *echo.Echo) {
	e.POST("/login", handlers.Login)
	e.POST("/forgotpassword", handlers.ForgotPassword)
	e.POST("/resetpassword", handlers.ResetPassword)
}
