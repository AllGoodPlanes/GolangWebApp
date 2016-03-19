package main

import (
	"fmt"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if user == "OK" && pass == "match" {
			fmt.Println(w, "Authorised")
			next.ServeHTTP(w, req)
		} else {
			fmt.Println(w, "Not Authorised")
			http.Redirect(w, req, "/signin/", http.StatusFound)
			return
		}
	})
}
