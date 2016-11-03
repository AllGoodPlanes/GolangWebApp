package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var user string
var pass string

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

func main() {

	mux := http.NewServeMux()

	fmt.Println("Listening...")
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/about/", About)

	loggedin := http.HandlerFunc(Signedin)
	mux.Handle("/signin/", Signin(loggedin))
	membernews := http.HandlerFunc(News)
	mux.Handle("/auth/membernews/", Auth(membernews))
	suggestions := http.HandlerFunc(Display)
	mux.Handle("/auth/suggestions/", Auth(suggestions))
	addsuggestions := http.HandlerFunc(Addcomment)
	mux.Handle("/auth/addcomment/", Auth(addsuggestions))
	mux.HandleFunc("/register/", Register)
	mux.HandleFunc("/verify/", Verify)
	mux.HandleFunc(STATIC_URL, StaticHandler)
	http.ListenAndServe(GetPort(), mux)

}

type Context struct {
	Title  string
	Static string
	User   string
}

func Home(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	render(w, "home", context)
}

func About(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About us"}
	render(w, "about", context)
}

func Signedin(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome member"}
	render(w, "signedin", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
	tmpl_list := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Info: No port detected in the environment, defaulting to :" + port)

	return ":" + port
}
