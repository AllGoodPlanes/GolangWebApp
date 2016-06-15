package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("GoWebAppCookie")
		if err != nil {
			fmt.Println("error")
			http.Redirect(w, req, "/signin/", http.StatusFound)
		} else {

			var result GolangWebAppSession
			sessionCopy1 := mongoSession.Copy()
			collection1 := sessionCopy1.DB(TestDatabase).C("GolangWebAppSession")
			err = collection1.Find(bson.M{"ID": cookie.Value}).One(&result)

			fmt.Println("GOlangWebAppSession", result.Expires)
			fmt.Println("GOlangWebAppSession", result.ID)
			fmt.Println("GOlangWebAppSession", result.Username)

			now := time.Now()
			cookietime := cookie.Expires
			diff := now.Sub(cookietime)

			if cookie.Value == result.ID && diff > 0 {

				fmt.Println(w, "Authorised")
				next.ServeHTTP(w, req)
			} else {
				fmt.Println(w, "Not Authorised")
				http.Redirect(w, req, "/signin/", http.StatusFound)
				return
			}

		}
	})
}
