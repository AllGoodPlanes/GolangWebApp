package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

type GolangWebAppSession struct {
	Username string    `bson:"Username"`
	ID       string    `bson:"ID"`
	Expires  time.Time `bson:"Expires"`
}

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

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
				errorloginTemplate.Execute(w, "Sorry... it looks like either your Password, or Username, is incorrect.")
				return
			}
			db.QueryRow("SELECT Verified FROM UserContact WHERE username = $1", username).Scan(&Verified)

			if Verified == "Yes" {
				match := bcrypt.CompareHashAndPassword([]byte(Password), password1)

				if match == nil {

					log.Println("Executing cookie & session")
					exp := time.Now().Add(30 * time.Minute)
					session := &Session{
						ID:     sessID(20),
						Expiry: exp,
					}

					cookie := http.Cookie{Name: "GoWebAppCookie", Value: session.ID, Domain: "golangwebapp.herokuapp.com", Path: "/auth/", Expires: session.Expiry}
					http.SetCookie(w, &cookie)

					fmt.Println("method:", req.Method)

					sessionCopy1 := mongoSession.Copy()
					defer sessionCopy1.Close()
					collection1 := sessionCopy1.DB(TestDatabase).C("GolangWebAppSession")

					err = collection1.Insert(&GolangWebAppSession{Username: username, ID: session.ID, Expires: session.Expiry})

					if err != nil {
						log.Printf("error is%v", &err)
						panic(err)
					}

					log.Println("Executing Auth")
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
