package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"git.npcompete.com/OSPSM_Servers/src/api/encrypt_decrypt"
	"git.npcompete.com/OSPSM_Servers/src/dbconnect"
	"git.npcompete.com/OSPSM_Servers/src/view"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func Login(c echo.Context) (err error) {
	frmt := new(view.Users)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}
	emailid := strings.TrimSpace(frmt.Email)
	password := strings.TrimSpace(frmt.Password)
	if len(emailid) > 0 && len(password) > 0 {
		dbconn := dbconnect.DB
		userview := view.Users{}
		dbconn.Where("email = ?", emailid).Find(&userview)
		fmt.Println("password check")

		// Comparing the password with the hash
		isCorrectPassword := checkPassword(password, userview.Password)
		if !isCorrectPassword {
			return c.JSON(http.StatusForbidden, view.GetFailureResponse("Password is incorrect", 0, nil))
		}

		// check username and password DB after hashing the password
		if emailid == userview.Email && isCorrectPassword && userview.Role == frmt.Role {
			// cookie := &http.Cookie{}

			// cookie.Name = "OSMPSM"
			// cookie.Value = "private"
			// cookie.Expires = time.Now().Add(15 * time.Minute)

			// c.SetCookie(cookie)

			// create jwt token
			token, err := createJwtToken(userview.Email)
			if err != nil {
				log.Println("Error Creating JWT token", err)
				return c.JSON(http.StatusInternalServerError, view.GetFailureResponse("something went wrong", 0, nil))
			}

			return c.JSON(http.StatusOK, view.GetSuccessResponse("user logged in successfully", map[string]string{
				"token": token,
				"role":  userview.Role}))
		}
	}

	return c.JSON(http.StatusUnauthorized, view.GetFailureResponse("Your username or password were wrong", 0, nil))
}

func checkPassword(password string, userpassword string) bool {
	key := "b826f19a79c492a6821233fdf54c5e923a5b839a357b024c62726821a02ba9f9"
	decrypt_password := encrypt_decrypt.DecryptData(userpassword, key)
	if decrypt_password != password {
		log.Println(decrypt_password)
		return false
	}
	return true
}

func createJwtToken(email string) (string, error) {
	claims := JwtClaims{
		email,
		jwt.StandardClaims{
			Id:        "ospsm_user",
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("gouthamospsm"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ForgotPassword(c echo.Context) (err error) {
	email := c.QueryParam("email")
	dbconn := dbconnect.DB
	userview := view.Users{}
	clientview := view.Client{}
	dbconn.Where("email = ?", email).Find(&clientview)
	if email != clientview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}
	dbconn.Where("email = ?", email).Find(&userview)

	key := "b826f19a79c492a6821233fdf54c5e923a5b839a357b024c62726821a02ba9f9"
	decrypt_password := encrypt_decrypt.DecryptData(userview.Password, key)

	ForgotSendEmail(clientview.ClientName, clientview.Email, decrypt_password)

	return c.JSON(http.StatusUnauthorized, view.GetSuccessResponse("password sent to your mail", nil))

}

func ResetPassword(c echo.Context) (err error) {
	frmt := new(view.Resetpass)
	if e := c.Bind(frmt); e != nil {
		err = c.JSON(http.StatusBadRequest, view.GetFailureResponse("Cannot process request", 0, nil))
		return
	}
	dbconn := dbconnect.DB
	userview := view.Users{}
	dbconn.Where("email = ?", frmt.Email).Find(&userview)
	if frmt.Email != userview.Email {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("User not Exists", 0, nil))
	}

	key := "b826f19a79c492a6821233fdf54c5e923a5b839a357b024c62726821a02ba9f9"
	decrypt_password := encrypt_decrypt.DecryptData(userview.Password, key)

	if frmt.OldPassword != decrypt_password {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("Old Password wrong", 0, nil))
	}

	if frmt.NewPassword != frmt.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, view.GetFailureResponse("New password and confirm password not match", 0, nil))
	}

	encryptpassword := encrypt_decrypt.EncryptData(frmt.ConfirmPassword, key)

	dbconn.Model(&userview).Where("email = ?", frmt.Email).Update("password", encryptpassword)

	return c.JSON(http.StatusUnauthorized, view.GetSuccessResponse("Password Changed Successfully", nil))

}
