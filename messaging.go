package main

import (
	"log"
	"net/smtp"
	"net/url"
	"fmt"
)

func emaillnk(e, u, h string) {
	// Set up authentication information.i
	fmt.Println(e)
	fmt.Println(u)
	fmt.Println(h)
	auth := smtp.PlainAuth("","richardlong@gocloudcoding.com", "2N!tefight", "smtp.office365.com")
	baseUrl := "http://golangwebapp.herokuapp.com/verify/?"
	params := url.Values{}
	params.Add("token", h)
	params.Add("email", e)
	finalUrl := baseUrl + params.Encode()

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{e}
	msg := []byte("To:" + e + "\r\n" +
		"Subject: Complete App Registration\r\n" +
		"\r\n" +
		"Click on link to complete App registration.\r\n" +
		finalUrl)
	err := smtp.SendMail("smtp.office365.com:587", auth, "richardlong@gocloudcoding.com", to, msg)

	if err != nil {
		log.Fatal(err)
	}
}
