package handlers

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"
)

func SendEmail(ClientName string, Email string, randomNum int) (sented bool, err error) {
	// sender data
	from := "gouthamotp@gmail.com"
	password := "7904446715"
	// receiver address
	toEmail := Email
	to := []string{toEmail}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := "Subject: OSPSM OTP\n"
	convert := strconv.Itoa(randomNum)
	body := "Hi\t" + ClientName +
		"\nsuccessfully created your account please find the otp and register using the link \t" + convert +
		"\nValid 10 Mins"
	message := []byte(subject + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// send mail
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	log.Println("mail communicated")
	e := smtp.SendMail(address, auth, from, to, message)

	if e != nil {
		fmt.Printf("smtp error: %s", e)
		return
	} else {
		sented = true
		log.Println("mail sent sucessfully")
	}

	return
}

func ForgotSendEmail(ClientName string, Email string, Password string) (sented bool, err error) {
	// sender data
	from := "gouthamotp@gmail.com"
	password := "7904446715"
	// receiver address
	toEmail := Email
	to := []string{toEmail}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := "Subject: OSPSM Forgot Password\n"
	body := "Hi\t" + ClientName +
		"\nYour Password is  \t" + Password
	message := []byte(subject + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// send mail
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	log.Println("mail communicated")
	e := smtp.SendMail(address, auth, from, to, message)

	if e != nil {
		fmt.Printf("smtp error: %s", e)
		return
	} else {
		sented = true
		log.Println("mail sent sucessfully")
	}

	return
}
