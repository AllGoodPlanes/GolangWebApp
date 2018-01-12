package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"compress/gzip"
	"strings"
)

var user string
var pass string

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

func main() {

	mux := http.NewServeMux()

	fmt.Println("Listening...")
	mux.HandleFunc("/", makeGzipHandler(Home))
	mux.HandleFunc("/about/", makeGzipHandler(About))
	mux.HandleFunc("/feedback/", makeGzipHandler(VisitorFbck))
	mux.HandleFunc("/signedin/",Signedin)
	loggedin := makeGzipHandler(http.HandlerFunc(Signedin))
	mux.Handle("/signin/", Signin(loggedin))
	membernews := makeGzipHandler(http.HandlerFunc(News))
	mux.Handle("/auth/membernews/", Auth(membernews))
	loggedout:= makeGzipHandler(http.HandlerFunc(Signout))
	mux.Handle("/auth/signout/", Auth(loggedout))
	suggestions := makeGzipHandler(http.HandlerFunc(Comments))
	mux.Handle("/auth/suggestions/", Auth(suggestions))
	addsuggestions := makeGzipHandler(http.HandlerFunc(Addcomment))
	mux.Handle("/auth/addcomment/", Auth(addsuggestions))
	vote := makeGzipHandler(http.HandlerFunc(V))
	mux.Handle("/auth/votes/",Auth(vote))
	mux.HandleFunc("/register/", makeGzipHandler(Register))
	mux.HandleFunc("/verify/", makeGzipHandler(Verify))
	fs:= http.FileServer(http.Dir("static"))
	mux.Handle("/static/",http.StripPrefix("/static/", fs))
//	img:= http.FileServer(http.Dir("images"))
//	mux.Handle("/image/",http.StripPrefix("/images", img))
	http.ListenAndServe(GetPort(), mux)


}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}




type Context struct {
	Title  string
	Static string
	User   string


}

func (w gzipResponseWriter) Write(b []byte) (int, error){
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r*http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"),"gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr:=gzipResponseWriter{Writer:gz, ResponseWriter: w}
		fn(gzr, r)
	}
}


func Home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "Welcome!"}
	render(w, "home", "Gopher", context)
}

func About(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "About us"}
	render(w, "about", "", context)
}

func render(w http.ResponseWriter, tmpl string, img string, context Context) {
	if img == ""{
	context.Static = STATIC_URL
	tmpl_list := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl),
	}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	//err = t.Execute(w, context )
	err= t.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}else{

	context.Static = STATIC_URL
	tmpl_list := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl),
		fmt.Sprintf("images/%s.png", img),
	}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	//err = t.Execute(w, context )
	err= t.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Print("template executing error: ", err)
}}

}


func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Info: No port detected in the environment, defaulting to :" + port)

	return ":" + port
}
