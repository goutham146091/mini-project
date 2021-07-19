package main

import (
	"fmt"

	"git.npcompete.com/OSPSM_Servers/src/dbconnect"
	"git.npcompete.com/OSPSM_Servers/src/router"
	_ "github.com/rizalgowandy/git.npcompete.com/OSPSM_Servers/src/main/docs" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	fmt.Println("Welcome to the webserver")
	dbconn, err := dbconnect.DbConnection()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("connection successfull")
	}
	e := router.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Start(":8000")
	fmt.Printf("Listen And Serve")
	defer dbconn.Close()
}
