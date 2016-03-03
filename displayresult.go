package main

import (
	"gopkg.in/mgo.v2/bson"
	"html/template"
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

func lookup(w http.ResponseWriter, req *http.Request) {

	memberareaTemplate.ExecuteTemplate(w, "memberarea.html", nil)

}
func display(w http.ResponseWriter, r *http.Request) {
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(TestDatabase).C("AddressData")

	var result AddressData

	err = collection.Find(bson.M{"Name": r.FormValue("name")}).One(&result)

	if result.Email != "" {
		context := Put{Result: "The e.mail address is:" + result.Email}
		displayresultTemplate.ExecuteTemplate(w, "displayresult.html", context)
	} else {
		context := Put{Result: "Sorry, that name isn't registered. Try again?"}
		displayresultTemplate.ExecuteTemplate(w, "displayresult.html", context)

	}

}
