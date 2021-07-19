package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"git.npcompete.com/OSPSM_Servers/src/api/encrypt_decrypt"
	"git.npcompete.com/OSPSM_Servers/src/dbconnect"
	"git.npcompete.com/OSPSM_Servers/src/view"
	"github.com/labstack/echo"
)

func RegisterBranch(c echo.Context) (err error) {
	frmt := new(view.EmailPassword)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB
	branchview := view.Branch{}
	dbconn.Where("email = ?", frmt.Email).Find(&branchview)

	if frmt.Email != branchview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}

	if frmt.Role != branchview.Role {
		return c.JSON(http.StatusUnauthorized, view.GetFailureResponse("unauthorized user", 0, nil))
	}

	if frmt.Password != frmt.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("password and confirm password not match", 0, nil))
	}

	if frmt.InviteCode != branchview.InviteCode {
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
			Role:           "Branch",
			CreateDateTime: formatedTime}
		dbconn.Create(&userview)
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}
	dbconn.Model(&branchview).Where("email = ?", frmt.Email).Update("is_registerd", true)
	return c.JSON(http.StatusOK, view.GetSuccessResponse("user registered successfully", map[string]string{
		"token": token,
		"role":  branchview.Role}))
}

func UpdateBranch(c echo.Context) (err error) {
	frmt := new(view.BranchDataPojo)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB
	branchview := view.Branch{}
	dbconn.Where("email = ?", frmt.Branch.Email).Find(&branchview)

	if frmt.Branch.Role != branchview.Role {
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
		UserID:         branchview.UserID,
		NumberType:     frmt.ContactNumber.NumberType,
		ISD_Code:       frmt.ContactNumber.ISD_Code,
		ContactNum:     frmt.ContactNumber.ContactNum,
		CreateDateTime: formatedTime,
	}
	dbconn.Create(&contactview)

	contactaddressview := view.ContactAddress{
		UserID:         branchview.UserID,
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

func DeleteBranch(c echo.Context) (err error) {
	email := c.QueryParam("email")

	dbconn := dbconnect.DB
	branchview := view.Branch{}
	dbconn.Where("email = ?", email).Find(&branchview)
	if email != branchview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}

	dbconn.Model(&view.Branch{}).Where("email = ?", email).Update("is_deleted", true)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("user deleted successfully", nil))
}

func DeleteProduct(c echo.Context) (err error) {
	pid := c.QueryParam("id")
	productid, err := strconv.Atoi(pid)

	if err != nil {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("Conversion Error", 0, nil))
	}

	dbconn := dbconnect.DB
	productview := view.Product{}
	dbconn.Where("productname = ?", pid).Find(&productview)
	if productid != productview.ProductID {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("Product not Exists", 0, nil))
	}

	dbconn.Where("productname = ?", pid).Delete(&productview)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("product deleted successfully", nil))
}

func GetProductList(c echo.Context) (err error) {
	role := c.QueryParam("role")
	email := c.QueryParam("email")
	branchview := view.Branch{}
	var productviews []view.Product

	dbconn := dbconnect.DB
	dbconn.Where("email = ?", email).Find(&branchview)
	if email != branchview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}

	if role != "Branch" {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("unauthorized", 0, nil))
	}

	dbconn.Find(&productviews)

	return c.JSON(http.StatusOK, view.GetSuccessResponse("client fetched successfully", productviews))
}

func GetBranchDetail(c echo.Context) (err error) {
	email := c.QueryParam("email")
	var branchview view.Branch
	var contactnumbers view.ContactNumber
	var contactAddress view.ContactAddress

	dbconn := dbconnect.DB

	err = dbconn.Find(&branchview, "email = ?", email).Error
	dbconn.Find(&contactnumbers, "user_id = ?", branchview.UserID)
	dbconn.Find(&contactAddress, "user_id = ?", branchview.UserID)

	result := view.UserDataPojo{
		Client: view.Client{
			UserID:         branchview.UserID,
			Email:          branchview.Email,
			ClientName:     branchview.BranchName,
			Role:           branchview.Role,
			CreateDateTime: branchview.CreateDateTime,
			InviteCode:     0,
			Is_Registered:  branchview.Is_Registered,
			Is_Deleted:     branchview.Is_Deleted,
			Is_Active:      branchview.Is_Active,
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

	return c.JSON(http.StatusOK, view.GetSuccessResponse("Branch Detail fetched successfully", result))
}

func AddProduct(c echo.Context) (err error) {
	frmt := new(view.Product)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB

	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)

	if len(frmt.ProductName) > 0 && frmt.Quantity > 0 && frmt.PricePerUnit > 0 &&
		frmt.WholePrice > 0 {
		productview := view.Product{
			ProductID:       0,
			BranchID:        0,
			ProductName:     frmt.ProductName,
			Quantity:        frmt.Quantity,
			PricePerUnit:    frmt.PricePerUnit,
			WholePrice:      frmt.WholePrice,
			CreateDateTime:  formatedTime,
			AvailableStatus: frmt.AvailableStatus,
		}
		dbconn.Create(&productview)
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}
	return c.JSON(http.StatusOK, view.GetSuccessResponse("user registered successfully", nil))
}

func UpdateProduct(c echo.Context) (err error) {
	frmt := new(view.Product)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB

	if len(frmt.ProductName) > 0 && frmt.Quantity > 0 && frmt.PricePerUnit > 0 &&
		frmt.WholePrice > 0 {

		dbconn.Model(&view.Branch{}).Where("product_name = ?", frmt.ProductName).Updates(
			map[string]interface{}{
				"productname":      frmt.ProductName,
				"quantity":         frmt.Quantity,
				"price_per_unit":   frmt.PricePerUnit,
				"whole_price":      frmt.WholePrice,
				"available_status": frmt.AvailableStatus,
			})
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}
	return c.JSON(http.StatusOK, view.GetSuccessResponse("user registered successfully", nil))
}

func MakeInvoice(c echo.Context) (err error) {
	frmt := new(view.Invoice)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}

	dbconn := dbconnect.DB

	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)

	if len(frmt.ProductName) > 0 && frmt.Quantity > 0 && frmt.PricePerUnit > 0 &&
		frmt.WholePrice > 0 && frmt.TotalPrice > 0 {
		total_price := frmt.Quantity * frmt.PricePerUnit
		invoiceview := view.Invoice{
			InvoiceID:    0,
			ProductID:    0,
			ProductName:  frmt.ProductName,
			Quantity:     frmt.Quantity,
			PricePerUnit: frmt.PricePerUnit,
			WholePrice:   frmt.WholePrice,
			PurchaseDate: formatedTime,
			Discount:     frmt.Discount,
			TotalPrice:   total_price,
		}
		dbconn.Create(&invoiceview)
	} else {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Bad request", 0, nil))
		return
	}
	return c.JSON(http.StatusOK, view.GetSuccessResponse("invoice created successfully", nil))
}
