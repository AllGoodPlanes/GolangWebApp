package main

import(

	"gopkg.in/mgo.v2/bson"
	"time"
	"net/http"
	"fmt"

)

type Feedback struct{
	Email string `bson:"Email"`
	Interest string `bson:"Interest"`
	Device string `bson:"Device"`
	Feedback string `bson:"Feedback"`
}

func VisitorFbck(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "Visitor Feedback"}
	fmt.Println("method:", req.Method)
	if req.Method == "GET"{
	render (w, "feedback", context)
} else {
	req.ParseForm()

	sessionCopy := mongoSession.Copy()
	collection := sessionCopy.DB(TestDatabase).C("Feedback")
	defer sessionCopy.Close()

	req.ParseForm()

	time := time.Now()

	email := req.FormValue("email")
	interest := req.FormValue("interest")
	device := req.FormValue("device")
	feedback := req.FormValue("feedback")

	newDoc := bson.M{"_id":time, "Email":email, "Interest":interest, "Device":device, "Feedback":feedback}

	collection.Insert(newDoc)

	http.Redirect(w, req,"/feedback/", http.StatusFound)

}
}



