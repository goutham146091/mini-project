package router

import (
	"git.npcompete.com/OSPSM_Servers/src/api"
	"git.npcompete.com/OSPSM_Servers/src/api/middlewares"
	"github.com/labstack/echo"
)

func New() *echo.Echo {
	api_url := echo.New()

	// create groups
	adminGroup := api_url.Group("/admin")
	clientGroup := api_url.Group("/client")

	// set all middlewares
	middlewares.SetApiMiddlewares(api_url)
	// middlewares.SetAdminMiddlewares(adminGroup)

	// set common routes
	api.CommonGroup(api_url)

	// set group routes
	api.AdminGroup(adminGroup)
	api.ClientGroup(clientGroup)

	return api_url
}
