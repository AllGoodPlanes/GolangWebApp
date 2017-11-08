package main

import(

"html/template"
"gopkg.in/mgo.v2/bson"
	"log"
"path/filepath"
"fmt"
"time"
	"net/http"
)

type Suggestions struct{
	Name string `bson:"Name"`
	Comment string `bson:"Comment"`
	Date string `bson:"Date"`
	Votes []Votes `bson:"Votes"`
	VoteTot int`bson:"VoteTot"`
	}

type Votes struct{
	Voter string `bson: "Voter"`
	Vote int `bson: "Vote"`

}

func Comments(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionCopy:= mongoSession.Copy()
	collection := sessionCopy.DB(TestDatabase).C("Suggestions")
	defer sessionCopy.Close()

	pipe:= collection.Pipe([]bson.M{

	bson.M{"$unwind": "$Votes"},
	bson.M{"$project":bson.M{"_id":1, "Comment":1, "Name":1, "Date":1, "C":bson.M{"$sum":1}}},
	bson.M{"$group":bson.M{"_id":"$Comment", "Name":bson.M{"$first":"$Name"}, "Date":bson.M{"$first":"$Date"},"C1":bson.M{"$sum":"$C"}}},
	bson.M{"$project":bson.M{"_id":1, "Name":1, "Date":1, "Count":bson.M{"$subtract":[]interface{}{"$C1",1}}}},
	bson.M{"$project":bson.M{"_id":1, "Name":1, "Date":1,"Count":1}},
	},
)

	var err2 error
	resp := []bson.M{}
	err2 = pipe.All(&resp)

	if err2 !=nil{
		panic(err2)
	}


	fmt.Println(resp)

	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "suggestions.html")
	dp := filepath.Join("templates", "display.html")

	t, err := template.ParseFiles(lp, fp, dp)
	if err!=nil{
		log.Fatalln(err)
	}


	err=t.ExecuteTemplate(w, "base", resp)
	if err!= nil {
		log.Fatalln(err)
	}

}

func Addcomment(w http.ResponseWriter,req *http.Request){
	w.Header().Set("Content-Type", "text/html")

	sessionCopy := mongoSession.Copy()

	cookie, err :=req.Cookie("GoWebAppCookie")
	if err !=nil {
		fmt.Println("cookie error")

	}
	time := time.Now()
	date := time.Format("2006-01-02 15:04")

	var result GolangWebAppSession
	collection := sessionCopy.DB(TestDatabase).C("GolangWebAppSession")
	err= collection.Find(bson.M{"ID":cookie.Value}).One(&result)
	collection1 := sessionCopy.DB(TestDatabase).C("Suggestions")
	defer sessionCopy.Close()
	req.ParseForm()
	name := result.Username
	suggestion:=req.FormValue("suggestion")

	newDoc := bson.M{"_id":time,"Comment": suggestion, "Name": name, "Date":date}

	collection1.Insert(newDoc)

	comment := bson.M{"Comment":suggestion}
	addvoter:= bson.M{"$addToSet":bson.M{"Votes":bson.M{"Voter":"K"}}}

	collection1.Update(comment, addvoter)

	var err1 error


	pipe:= collection1.Pipe([]bson.M{

	bson.M{"$unwind": "$Votes"},
	bson.M{"$project":bson.M{"_id":1, "Comment":1, "Name":1, "Date":1, "C":bson.M{"$sum":1}}},
	bson.M{"$group":bson.M{"_id":"$Comment", "Name":bson.M{"$first":"$Name"}, "Date":bson.M{"$first":"$Date"},"C1":bson.M{"$sum":"$C"}}},
	bson.M{"$project":bson.M{"_id":1, "Name":1, "Date":1, "Count":bson.M{"$subtract":[]interface{}{"$C1",1}}}},
	bson.M{"$project":bson.M{"_id":1, "Name":1, "Date":1,"Count":1}},

	},
)

	var err2 error
	resp := []bson.M{}
	err2 = pipe.All(&resp)

	if err2 !=nil{
		panic(err2)
	}
	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "suggestions.html")
	dp := filepath.Join("templates", "display.html")

	t, err1 := template.ParseFiles(lp, fp, dp)
	if err1!=nil{
		log.Fatalln(err1)
	}


	err1=t.ExecuteTemplate(w, "base", resp)
	if err1!= nil {
		log.Fatalln(err1)
	}
}

func V(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/html")

	sessionCopy := mongoSession.Copy()

	cookie, err :=req.Cookie("GoWebAppCookie")
	if err !=nil {
		fmt.Println("cookie error")

	}

	var result GolangWebAppSession
	collection := sessionCopy.DB(TestDatabase).C("GolangWebAppSession")
	err= collection.Find(bson.M{"ID":cookie.Value}).One(&result)
	collection1 := sessionCopy.DB(TestDatabase).C("Suggestions")
	defer sessionCopy.Close()
	req.ParseForm()

	name := result.Username
	v := req.FormValue("vote")
	comment := bson.M{"Comment":v}
	addvote:= bson.M{"$addToSet":bson.M{"Votes":bson.M{"Voter":name}}}

	err = collection1.Update(comment, addvote,)

	pipe:= collection1.Pipe([]bson.M{

	bson.M{"$unwind": "$Votes"},
	bson.M{"$project":bson.M{"_id":1, "Comment":1, "Name":1, "Date":1, "C":bson.M{"$sum":1}}},
	bson.M{"$group":bson.M{"_id":"$Comment", "Name":bson.M{"$first":"$Name"}, "Date":bson.M{"$first":"$Date"},"C1":bson.M{"$sum":"$C"}}},
	bson.M{"$project":bson.M{"_id":1, "Name":1, "Date":1, "Count":bson.M{"$subtract":[]interface{}{"$C1",1}}}},
	bson.M{"$project":bson.M{"_id":1, "Name":1, "Date":1,"Count":1}},

	},
)

	var err2 error
	resp := []bson.M{}
	err2 = pipe.All(&resp)

	if err2 !=nil{
		panic(err2)
	}

	fmt.Println(resp)


	var err1 error


	lp := filepath.Join("templates", "base.html")
	fp := filepath.Join("templates", "suggestions.html")
	dp := filepath.Join("templates", "display.html")

	t, err1 := template.ParseFiles(lp, fp, dp)
	if err1!=nil{
		log.Fatalln(err1)
	}


	err1=t.ExecuteTemplate(w, "base", resp)
	if err1!= nil {
		log.Fatalln(err1)
	}

	}
