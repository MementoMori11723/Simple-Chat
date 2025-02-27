package Server

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"simple-chat/server/functions"
	"simple-chat/server/middleware"
)

var (
	//go:embed pages/*html
	files  embed.FS
	pages  = "pages/"
	layout = pages + "layout.html"

	indexTmpl *template.Template
	errorTmpl *template.Template
	aboutTmpl *template.Template
	loginTmpl *template.Template

	routes = map[string]http.HandlerFunc{
		"/":      home,
		"/about": about,
		"/login": login,
		"/error": errorFunc,
		"/404":   pageNotFound,

		"POST /signup": functions.Signup,
		"POST /login":  functions.Login,
	}
)

type _error struct {
	ErrorCode  int
	ErrorTitle string
	ErrorMsg   string
}

func init() {
	indexTmpl = getTemplate("home")
	errorTmpl = getTemplate("error")
	aboutTmpl = getTemplate("about")
	loginTmpl = getTemplate("login")
	go functions.TestConnection()
}

func getTemplate(name string) *template.Template {
	tmpl, err := template.ParseFS(
		files, layout, fmt.Sprintf("%s%s.html", pages, name),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return tmpl
}

func New() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/ws/{id}", functions.SocketHandler())
	for route, handler := range routes {
		mux.Handle(route, middleware.Logger(handler))
	}
	return mux
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
	if err := indexTmpl.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	if err := aboutTmpl.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if err := loginTmpl.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	if err := errorTmpl.Execute(w, _error{
		ErrorCode:  http.StatusNotFound,
		ErrorTitle: "Page Not Found",
		ErrorMsg:   "Sorry! the page you are looking for is either moved or does not exist.",
	}); err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
}

func errorFunc(w http.ResponseWriter, _ *http.Request) {
	if err := errorTmpl.Execute(w, _error{
		ErrorCode:  http.StatusInternalServerError,
		ErrorTitle: "Internal Server Error",
		ErrorMsg:   "Sorry! something went wrong on our end and we are working to fix it, please try again later.",
	}); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
