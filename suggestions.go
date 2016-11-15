package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

type Suggestions struct{
	Username string `bson:"Name"`
	Comment string `bson:"Comment"`
}

type Put struct{
	Restlt string
	Static string
	User string
}


var displaysuggestionsTemplate = template.Must(template.ParseGlob("templates/suggestions.html"))
var err error

func Addcomment(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	sessionCopy:= mongoSession.Copy()

	cookie, err := req.Cookie("GoWebAppCookie")
	if err !=nil {
		fmt.Println("cookie error")
	}

	var result1 GolangWebAppSession
	collection1 := sessionCopy.DB(TestDatabase).C("GolangWebAppSession")
	err = collection1.Find(bson.M{"ID": cookie.Value}).One(&result1)
	collection2 := sessionCopy.DB(TestDatabase).C("Suggestions")
	defer sessionCopy.Close()
	req.ParseForm()
	name := result1.Username
	suggestion := req.FormValue("suggestion")
	err =collection2.Insert(&Suggestions{Comment:suggestion,Username: name})
}

func Display(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	displaysuggestionsTemplate.ExecuteTemplate(w, "suggestions.html", nil)
}

