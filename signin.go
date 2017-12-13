package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
	"strings"
	"path/filepath"
)

type Session struct {
	ID     string
	UserID string
	Expiry string
}

type GolangWebAppSession struct {
	Username string    `bson:"Username"`
	ID       string    `bson:"ID"`
	//Expires  time.Time `bson:"Expires"`
	Expires string `bson:"Expires"`

}

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
//	Expires    time.Time
	Expires string
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

func Signedin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Context-Type", "text/html")

	url := req.URL.String()
	s := strings.Split(url,"/")
	ip, port := s[1], s[2]
	fmt.Println(ip)
	fmt.Println(port)

	sessionCopy:= mongoSession.Copy()
	collection := sessionCopy.DB(TestDatabase).C("GolangWebAppSession")
	defer sessionCopy.Close()

	X := time.Now()
	NTm1 := X.Add(3 * time.Minute)
	NTm := NTm1.Format("2006-01-02 15:04")

	pipe:= collection.Pipe([]bson.M{
	bson.M{"$match":bson.M{"Username":port}},
	bson.M{"$group":bson.M{"_id":"$_id", "Expires":bson.M{"$push":"$Expires"}}},
	bson.M{"$match":bson.M{"Expires":bson.M{"$ne":NTm}}},
	bson.M{"$project":bson.M{"_id":1,"Expires":1}},
	bson.M{"$sort":bson.M{"Expires":1}},
	bson.M{"$group":bson.M{"_id":0, "Expires":bson.M{"$last":"$Expires"}}},
	bson.M{"$project":bson.M{"Expires":1}},
	},
)

	var err2 error
	resp :=[]bson.M{}
	err2 = pipe.All(&resp)

	if err2 !=nil{
		panic(err2)
	}

	a := fmt.Sprintf("%s",resp)
	b:= strings.Replace(a, "[map[_id:%!s(int=0) Expires:[","", -1)
	c := strings.Replace(b, "]]]", "", -1)
	d := fmt.Sprintf("%s", c)

	var f string
	if d == "[]"{
		f = "new member"
	}else{
		f = d
	}

	e := fmt.Sprintf("%s", port)

	type Z struct{
	Name string
	Expires string
	}

	var x =[]Z{{e,f}}
	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "signedin.html")

	t, err := template.ParseFiles(lp, fp)
	if err!=nil{
		log.Fatalln(err)
	}


	err=t.ExecuteTemplate(w, "base", x)
	if err!= nil {
		log.Fatalln(err)

	}
}

func Signout(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "You've logged out successfully"}
	render(w, "signout", context)

	cookie, err := req.Cookie("GoWebAppCookie")
	if err != nil {
		fmt.Println("error")
		http.Redirect(w,req, "/signin/", http.StatusFound)
	}
	sessionCopy := mongoSession.Copy()
	collection := sessionCopy.DB(TestDatabase).C("GolangWebAppSession")
	defer sessionCopy.Close()

	fmt.Println(cookie.Value)


	Loggedin := bson.M{"ID":cookie.Value}
	IDamended := bson.M{"$set":bson.M{"ID":"xxx"}}

	collection.Update(Loggedin, IDamended)


}

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
					exp1 := time.Now().Add(3 * time.Minute)
					exp:= exp1.Format("2006-01-02 15:04")

					session := &Session{
						ID:     sessID(20),
						Expiry: exp,
					}

					cookie := http.Cookie{Name: "GoWebAppCookie", Value: session.ID, Domain: "golangwebapp.herokuapp.com", Path: "/auth/", Expires: exp1}
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

				http.Redirect(w, req, "/signedin/"+username, http.StatusFound)

				return
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


