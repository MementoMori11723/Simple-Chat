package Server

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"simple-chat/server/functions"
	"simple-chat/server/middleware"
)

var (
	//go:embed pages/*html
	files  embed.FS
	pages  = "pages/"
	layout = pages + "layouts.html"

	indexTmpl *template.Template
	errorTmpl *template.Template

	routes = map[string]http.HandlerFunc{
		"/":      home,
		"/404":   pageNotFound,
		"/error": errorFunc,
	}
)

func init() {
	var err error
	indexTmpl, err = template.ParseFS(
		files, layout, pages+"index.html",
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	errorTmpl, err = template.ParseFS(
		files, layout, pages+"error.html",
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}
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

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	if err := errorTmpl.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
}

func errorFunc(w http.ResponseWriter, _ *http.Request) {
	if err := errorTmpl.Execute(w, nil); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
