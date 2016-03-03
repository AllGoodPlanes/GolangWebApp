package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v1"
	"html/template"

	"net/http"
	"regexp"
)

type UserContact struct {
	ID       int64
	Verified string `db:"Verified"`
	Username string `db:"UserName"`
	EMail    string `db:"EMail"`
	Password string `db:"Password"`
	Country  string `db:"Country"`
	City     string `db:"City"`
	PostCode string `db:"PostCode"`
	Mobile   string `db:"Mobile"`
}

var emailregnotificationTemplate = template.Must(template.ParseGlob("templates/regnotification.html"))
var registererrorsTemplate = template.Must(template.ParseGlob("templates/registererrors.html"))

func Verify(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET params:", req.URL.Query())
	email := req.URL.Query().Get("email")
	token := req.URL.Query().Get("token")
	if token != "" {
		fmt.Println(token)
		fmt.Println(email)

		var Verified string

		match := bcrypt.CompareHashAndPassword([]byte(token), []byte(email))
		fmt.Println(match)

		if match == nil {
			db.QueryRow("UPDATE UserContact SET Verified = 'Yes' WHERE email = $1 RETURNING Verified", email).Scan(&Verified)
			memberareaTemplate.ExecuteTemplate(w, "memberarea.html", nil)
		} else {
			errorloginTemplate.Execute(w, "Sorry... it looks like something went wrong. Try the e.mail link again., then reregistering if there are still problems.")
		}
	}
}

func Register(w http.ResponseWriter, req *http.Request) {

	context := Context{Title: "Register new member"}
	fmt.Println("method:", req.Method) //get request method
	if req.Method == "GET" {
		render(w, "register", context)
	} else {
		req.ParseForm()
		// logic part of log in
		verified := "No"
		username := req.PostFormValue("username")
		email := req.PostFormValue("email")
		password := req.PostFormValue("password")
		reinputpswd := req.PostFormValue("reinputpassword")
		country := req.PostFormValue("country")
		city := req.PostFormValue("city")
		postCode := req.PostFormValue("postcode")
		mobile := req.PostFormValue("mobileno")

		var err error

		dialect := gorp.PostgresDialect{}
		dbmap := &gorp.DbMap{Db: db, Dialect: dialect}
		dbmap.DropTables()
		dbmap.AddTable(UserContact{}).SetKeys(true, "ID")
		err = dbmap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		err = dbmap.Db.Ping()
		if err != nil {
			panic(err)
		}
		password1 := []byte(password)
		email1 := []byte(email)

		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password1, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hashedPassword))

		hashedemail, err := bcrypt.GenerateFromPassword(email1, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hashedemail))

		re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
		var Username string
		var Email string
		passwordMin = 5
		passwordMax = 20
		paswd := password

		err = db.QueryRow("SELECT Username FROM UserContact WHERE username = $1", username).Scan(&Username)
		err = db.QueryRow("SELECT EMail FROM UserContact WHERE email = $1", email).Scan(&Email)

		if username == "" {
			registererrorsTemplate.Execute(w, "You must include a username.")
		} else if !re.MatchString(username) {
			registererrorsTemplate.Execute(w, "username must contain only alpha numeric & underscore characters.")
		} else if Username == username {
			registererrorsTemplate.Execute(w, "username already exists.")
		} else if email == "" {
			registererrorsTemplate.Execute(w, "You must include a valid e.mail address.")
		} else if Email == email {
			registererrorsTemplate.Execute(w, "That e.mail address is already registered.")
		} else if !checkLength(password, passwordMin, passwordMax) {
			registererrorsTemplate.Execute(w, "You must include a password with more than 5 characters.")
		} else if !passwordCk(paswd) {
			registererrorsTemplate.Execute(w, "Password must include UPPER & lower case letters and numbers.")
		} else if password != reinputpswd {
			registererrorsTemplate.Execute(w, "The re-input password doesn't match.")
		} else {
			emaillnk(email, username, string(hashedemail))
			c1 := &UserContact{0, verified, username, email, string(hashedPassword), country, city, postCode, mobile}
			dbmap.Insert(c1)
			emailregnotificationTemplate.ExecuteTemplate(w, "regnotification.html", nil)
		}
		fmt.Println("GET params:", req.URL.Query())
		token := req.URL.Query()["token"]
		if token != nil {
			fmt.Println(token)
		}
	}

}
