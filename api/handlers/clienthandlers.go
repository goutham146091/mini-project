package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"git.npcompete.com/OSPSM_Servers/src/api/encrypt_decrypt"
	"git.npcompete.com/OSPSM_Servers/src/dbconnect"
	"git.npcompete.com/OSPSM_Servers/src/view"
	"github.com/labstack/echo"
)

func RegisterClient(c echo.Context) (err error) {
	frmt := new(view.EmailPassword)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB
	clientview := view.Client{}
	dbconn.Where("email = ?", frmt.Email).Find(&clientview)

	if frmt.Email != clientview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}

	if frmt.Role != clientview.Role {
		return c.JSON(http.StatusUnauthorized, view.GetFailureResponse("unauthorized user", 0, nil))
	}

	if frmt.Password != frmt.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("password and confirm password not match", 0, nil))
	}

	if frmt.InviteCode != clientview.InviteCode {
		return c.JSON(http.StatusUnauthorized, view.GetFailureResponse("wrong invite code", 0, nil))
	}

	key := "b826f19a79c492a6821233fdf54c5e923a5b839a357b024c62726821a02ba9f9"

	encryptedPassword := encrypt_decrypt.EncryptData(frmt.Password, key)

	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)

	token, err := createJwtToken(frmt.Email)
	if err != nil {
		log.Println("Error Creating JWT token", err)
		return c.JSON(http.StatusInternalServerError, view.GetFailureResponse("something went wrong", 0, nil))
	}
	if len(frmt.Email) > 0 {
		userview := view.Users{
			Email:          frmt.Email,
			Password:       encryptedPassword,
			Role:           "Client",
			CreateDateTime: formatedTime}
		dbconn.Create(&userview)
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}
	dbconn.Model(&clientview).Where("email = ?", frmt.Email).Update("is_registered", true)
	return c.JSON(http.StatusOK, view.GetSuccessResponse("user registered successfully", map[string]string{
		"token": token,
		"role":  clientview.Role}))
}

func UpdateClient(c echo.Context) (err error) {
	frmt := new(view.UserDataPojo)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB
	clientview := view.Client{}
	dbconn.Where("email = ?", frmt.Client.Email).Find(&clientview)

	if frmt.Client.Role != clientview.Role {
		return c.JSON(http.StatusUnauthorized, view.GetFailureResponse("unauthorized user", 0, nil))
	}

	if len(frmt.ContactNumber.NumberType) < 0 || len(frmt.ContactNumber.ISD_Code) < 0 {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("numbertype or isd code empty", 0, nil))
	}

	if len(frmt.ContactAddress.City) < 0 || len(frmt.ContactAddress.State) < 0 || len(frmt.ContactAddress.Postalcode) < 0 {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("city or state or postalcode code empty", 0, nil))
	}

	if len(frmt.ContactAddress.AddressLine1) < 0 {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("city or state or postalcode code empty", 0, nil))
	}

	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)

	contactview := view.ContactNumber{
		UserID:         clientview.UserID,
		NumberType:     frmt.ContactNumber.NumberType,
		ISD_Code:       frmt.ContactNumber.ISD_Code,
		ContactNum:     frmt.ContactNumber.ContactNum,
		CreateDateTime: formatedTime,
	}
	dbconn.Create(&contactview)

	contactaddressview := view.ContactAddress{
		UserID:         clientview.UserID,
		AddressLine1:   frmt.ContactAddress.AddressLine1,
		AddressLine2:   frmt.ContactAddress.AddressLine2,
		City:           frmt.ContactAddress.City,
		State:          frmt.ContactAddress.State,
		Postalcode:     frmt.ContactAddress.Postalcode,
		CreateDateTime: formatedTime,
	}
	dbconn.Create(&contactaddressview)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("user data updated successfully", nil))
}

func DeleteClient(c echo.Context) (err error) {
	email := c.QueryParam("email")

	dbconn := dbconnect.DB
	clientview := view.Client{}
	dbconn.Where("email = ?", email).Find(&clientview)
	if email != clientview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}

	dbconn.Model(&view.Client{}).Where("email = ?", email).Update("is_deleted", true)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("user deleted successfully", nil))
}

// GetClientList godoc
// @Summary Show the list of the clients.
// @Description get the client list.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func GetClientList(c echo.Context) (err error) {
	role := c.QueryParam("role")
	var clientviews []view.Client

	dbconn := dbconnect.DB
	if role != "client" {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("unauthorized", 0, nil))
	}

	dbconn.Find(&clientviews)
	//dbconn.Select([]string{"client_name", "shoptype"}).Find(&clientviews)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("client fetched successfully", clientviews))
}

func GetClientDetail(c echo.Context) (err error) {
	email := c.QueryParam("email")
	var clientviews view.Client
	var contactnumbers view.ContactNumber
	var contactAddress view.ContactAddress

	dbconn := dbconnect.DB

	err = dbconn.Find(&clientviews, "email = ?", email).Error
	dbconn.Find(&contactnumbers, "user_id = ?", clientviews.UserID)
	dbconn.Find(&contactAddress, "user_id = ?", clientviews.UserID)

	result := view.UserDataPojo{
		Client: view.Client{
			UserID:         clientviews.UserID,
			Email:          clientviews.Email,
			ClientName:     clientviews.ClientName,
			Role:           clientviews.Role,
			CreateDateTime: clientviews.CreateDateTime,
			Shoptype:       clientviews.Shoptype,
			NoOfBranches:   clientviews.NoOfBranches,
			InviteCode:     0,
			Is_Registered:  clientviews.Is_Registered,
			Is_Deleted:     clientviews.Is_Deleted,
			Is_Active:      clientviews.Is_Active,
		},
		ContactNumber: view.ContactNumber{
			UserID:         0,
			NumberID:       0,
			NumberType:     contactnumbers.NumberType,
			ISD_Code:       contactnumbers.ISD_Code,
			ContactNum:     contactnumbers.ContactNum,
			CreateDateTime: contactnumbers.CreateDateTime,
		},
		ContactAddress: view.ContactAddress{
			UserID:         0,
			AddressID:      0,
			AddressLine1:   contactAddress.AddressLine1,
			AddressLine2:   contactAddress.AddressLine2,
			City:           contactAddress.City,
			State:          contactAddress.State,
			Postalcode:     contactAddress.Postalcode,
			CreateDateTime: contactAddress.CreateDateTime,
		},
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exist", 0, nil))
	}

	return c.JSON(http.StatusOK, view.GetSuccessResponse("client Detail fetched successfully", result))
}

func CreateBranch(c echo.Context) (err error) {
	email := c.QueryParam("email")
	frmt := new(view.Branch)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB
	branchview := view.Branch{}
	clientview := view.Client{}
	dbconn.Where("email = ?", frmt.Email).Find(&branchview)
	dbconn.Where("email = ?", email).Find(&clientview)

	var NoOfBranches int

	if frmt.Email == branchview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("Email Already Exists", 0, nil))
	}

	if NoOfBranches >= clientview.NoOfBranches {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("Client have Limited Branch only", 0, nil))
	}

	//create random numbers
	randomNum := random(1000, 9000)

	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)
	if len(frmt.Email) > 0 && len(frmt.Role) > 0 {
		clientview := view.Branch{
			BranchName:     frmt.BranchName,
			Email:          frmt.Email,
			BranchLocation: frmt.BranchLocation,
			Role:           frmt.Role,
			InviteCode:     randomNum,
			CreateDateTime: formatedTime}
		dbconn.Create(&clientview)
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}

	fmt.Println("branch registered")

	sent, err := SendEmail(frmt.BranchName, frmt.Email, randomNum)

	if !sent || err != nil {
		return c.JSON(http.StatusForbidden, view.GetFailureResponse("Mail not sent", 0, nil))
	} else {
		return c.JSON(http.StatusOK, view.GetSuccessResponse("user_branch Registered  successfully", nil))
	}

}

func GetBranchList(c echo.Context) (err error) {
	role := c.QueryParam("role")
	var branchviews []view.Branch

	dbconn := dbconnect.DB
	if role != "Branch" {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("unauthorized", 0, nil))
	}

	dbconn.Find(&branchviews)
	//dbconn.Select([]string{"client_name", "shoptype"}).Find(&clientviews)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("client fetched successfully", branchviews))
}
