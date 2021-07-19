package api

import (
	"git.npcompete.com/OSPSM_Servers/src/api/handlers"
	"github.com/labstack/echo"
)

func ClientGroup(adminGroup *echo.Group) {
	adminGroup.POST("/setpassword", handlers.RegisterClient)
	adminGroup.POST("/updateclient", handlers.UpdateClient)
	adminGroup.GET("/clientdetail", handlers.GetClientDetail)
	adminGroup.DELETE("/deleteclient", handlers.DeleteClient)
	adminGroup.POST("/registerbranch", handlers.CreateBranch)
	adminGroup.GET("/branchlist", handlers.GetBranchList)

}
