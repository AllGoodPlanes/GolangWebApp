package main

import (
	"fmt"
	"log"
	"net/http"
)

func News(w http.ResponseWriter, req *http.Request) {

	log.Println("Executing middlewareOne")

	fmt.Fprintf(w, "latest member news!")

}
