package main

import (
	"log"
	"net/smtp"
	"net/url"
)

func emaillnk(e, u, h string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "ricpelong@hotmail.co.uk", "2shoes2loose", "smtp.live.com")
	baseUrl := "http://golangwebapp.herokuapp.com/verify/?"
	params := url.Values{}
	params.Add("token", h)
	params.Add("email", e)
	finalUrl := baseUrl + params.Encode()

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{e}
	msg := []byte("To:" + u + "\r\n" +
		"Subject: Complete App Registration\r\n" +
		"\r\n" +
		"Click on link to complete App registration.\r\n" +
		finalUrl)
	err := smtp.SendMail("smtp.live.com:587", auth, "ricpelong@hotmail.co.uk", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
