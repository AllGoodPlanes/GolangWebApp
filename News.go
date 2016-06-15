package main

import (
	"net/http"
)

func News(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "All the latest news for members"}
	render(w, "membernews", context)

}
