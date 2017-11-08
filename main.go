package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
//	"time"
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

	loggedin := makeGzipHandler(http.HandlerFunc(Signedin))
	mux.Handle("/signin/", Signin(loggedin))
	membernews := makeGzipHandler(http.HandlerFunc(News))
	mux.Handle("/auth/membernews/", Auth(membernews))
	suggestions := makeGzipHandler(http.HandlerFunc(Comments))
	mux.Handle("/auth/suggestions/", Auth(suggestions))
	addsuggestions := makeGzipHandler(http.HandlerFunc(Addcomment))
	mux.Handle("/auth/addcomment/", Auth(addsuggestions))
	vote := makeGzipHandler(http.HandlerFunc(V))
	mux.Handle("/auth/votes/",Auth(vote))
	mux.HandleFunc("/register/", makeGzipHandler(Register))
	mux.HandleFunc("/verify/", makeGzipHandler(Verify))
	//mux.HandleFunc(STATIC_URL, makeGzipHandler(StaticHandler))
	fs:= http.FileServer(http.Dir("static"))
	mux.Handle("/static/",http.StripPrefix("/static/", fs))
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
	render(w, "home", context)
}

func About(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "About us"}
	render(w, "about",  context)
}

func Signedin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	context := Context{Title: "Welcome member"}
	render(w, "signedin", context, )
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
	tmpl_list := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	//err = t.Execute(w, context )
	err= t.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

//func StaticHandler(w http.ResponseWriter, req *http.Request) {
//	w.Header().Set("Content-Type","text/css")
//	static_file := req.URL.Path[len(STATIC_URL):]
//	if len(static_file) != 0 {
//		f, err := http.Dir(STATIC_ROOT).Open(static_file)
//		if err == nil {
//			content := io.ReadSeeker(f)
//			http.ServeContent(w, req, static_file, time.Now(), content)
//			return
//		}
//	}
//	http.NotFound(w, req)
//}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Info: No port detected in the environment, defaulting to :" + port)

	return ":" + port
}
