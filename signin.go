package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

var memberareaTemplate = template.Must(template.ParseGlob("templates/memberarea.html"))
var errorloginTemplate = template.Must(template.New("display").Parse(errorloginTemplateHTML))

func Signin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		context := Context{Title: "Member Sign In"}
		fmt.Println("method:", req.Method)
		if req.Method == "GET" {
			render(w, "signin", context)
		} else {
			req.ParseForm()
			username := req.PostFormValue("username")
			password := req.PostFormValue("password")

			password1 := []byte(password)
			var Password string
			var Verified string
			var err error

			err = db.QueryRow("SELECT Password FROM UserContact WHERE username = $1", username).Scan(&Password)

			if err != nil {
				log.Fatal(err)
			}
			db.QueryRow("SELECT Verified FROM UserContact WHERE username = $1", username).Scan(&Verified)

			if Verified == "Yes" {
				match := bcrypt.CompareHashAndPassword([]byte(Password), password1)

				if match == nil {

					next.ServeHTTP(w, req)
				} else {
					errorloginTemplate.Execute(w, "Sorry... it looks like either your Password, or Username, is incorrect.")
				}
			} else {
				errorloginTemplate.Execute(w, "Sorry... it looks like your e.mail address hasn't been verified. If you didn't get an e.mail message with a link, try registering again.")
			}
		}
	})
}

const errorloginTemplateHTML = `<!doctype html>
<html>
    <head>
<style>
  body {background-color:white}
  p    {color: orange}
legend {color: blue}
</style>
        <title></title>
    </head>
    <body>
        <form>
              <fieldset>
                        <legend></legend>
                        <p>
                        <p><b>{{html .}}</b></p>
                        <p><a href="/signin/">Try Loging in again!</a></p>
                        </p>
              </fielsset>
        </form>
    </body>
</html>`
