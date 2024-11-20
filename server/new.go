package Server

import (
	"embed"
	"net/http"
)

var (
	//go:embed pages/*html
	files  embed.FS
	layout = "pages/layouts.html"
)

func New() *http.ServeMux {
  mux := http.NewServeMux()
  mux.HandleFunc("/", homeHandler)
  mux.HandleFunc("/error", errorHandler)
  return mux
}
