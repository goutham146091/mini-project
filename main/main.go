package main

import (
	"fmt"

	"github.com/goutham146091/mini-project/tree/master/dbconnect"
	"github.com/goutham146091/mini-project/tree/master/router"
)

func main() {
	fmt.Println("Welcome to the webserver")
	dbconn, err := dbconnect.DbConnection()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("connection successfull")
	}
	e := router.New()
	e.Start(":8000")
	fmt.Printf("Listen And Serve")
	defer dbconn.Close()
}
