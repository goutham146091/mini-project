package handlers

import (
	"fmt"
	m "math/rand"
	"net/http"
	"time"

	"git.npcompete.com/OSPSM_Servers/src/dbconnect"
	"git.npcompete.com/OSPSM_Servers/src/view"
	"github.com/labstack/echo"
)

func random(min int, max int) int {
	return m.Intn(max-min) + min
}

func CreateClient(c echo.Context) (err error) {
	frmt := new(view.Client)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB
	clientview := view.Client{}
	dbconn.Where("email = ?", frmt.Email).Find(&clientview)

	if frmt.Email == clientview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("Email Already Exists", 0, nil))
	}

	//create random numbers
	randomNum := random(1000, 9000)

	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)
	if len(frmt.Email) > 0 && len(frmt.Role) > 0 && len(frmt.Shoptype) > 0 && frmt.NoOfBranches > 0 {
		clientview := view.Client{
			ClientName:     frmt.ClientName,
			Email:          frmt.Email,
			Role:           frmt.Role,
			Shoptype:       frmt.Shoptype,
			NoOfBranches:   frmt.NoOfBranches,
			InviteCode:     randomNum,
			CreateDateTime: formatedTime}
		dbconn.Create(&clientview)
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}

	fmt.Println("client registered")

	sent, err := SendEmail(frmt.ClientName, frmt.Email, randomNum)

	if !sent || err != nil {
		return c.JSON(http.StatusForbidden, view.GetFailureResponse("Mail not sent", 0, nil))
	} else {
		return c.JSON(http.StatusOK, view.GetSuccessResponse("user Registered  successfully", nil))
	}

}
