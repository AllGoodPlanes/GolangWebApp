package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

var memberareaTemplate = template.Must(template.ParseGlob("templates/memberarea.html"))

func Lookup(w http.ResponseWriter, req *http.Request) {

	fmt.Println("method:", req.Method)
	log.Println("Executing lookup")

	memberareaTemplate.ExecuteTemplate(w, "memberarea.html", nil)

}
