package api

import (
	"git.npcompete.com/OSPSM_Servers/src/api/handlers"
	"github.com/labstack/echo"
)

func AdminGroup(adminGroup *echo.Group) {
	adminGroup.POST("/register", handlers.CreateClient)
	adminGroup.GET("/clientlist", handlers.GetClientList)
}
