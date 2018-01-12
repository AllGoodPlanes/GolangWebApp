package main

import (
	"net/http"
)

func News(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "All the latest news for members"}
	render(w, "membernews", "", context)

}
