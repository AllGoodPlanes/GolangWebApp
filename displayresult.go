package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

type AddressData struct {
	Name  string `bson:"Name"`
	Email string `bson:"Email"`
}

type Put struct {
	Result string
	Static string
	User   string
}

var displayresultTemplate = template.Must(template.ParseGlob("templates/displayresult.html"))
var err error

func Display(w http.ResponseWriter, req *http.Request) {
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(TestDatabase).C("AddressData")

	var result AddressData

	err = collection.Find(bson.M{"Name": req.FormValue("name")}).One(&result)
	fmt.Println("method:", req.Method)
	log.Println("Executing display")
	if result.Email != "" {
		context := Put{Result: "The e.mail address is:" + result.Email}
		displayresultTemplate.ExecuteTemplate(w, "displayresult.html", context)
	} else {
		context := Put{Result: "Sorry, that name isn't registered. Try again?"}
		displayresultTemplate.ExecuteTemplate(w, "displayresult.html", context)

	}

}
