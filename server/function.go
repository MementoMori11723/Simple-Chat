package Server

import (
	"fmt"
	"html/template"
	"net/http"
)



func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	tmpl, err := template.ParseFS(files, layout, "pages/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error parsing template: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error executing template: ", err)
		return
	}
}

func errorHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFS(files, layout, "pages/error.html")
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("Error executing template: ", err)
		return
	}
}
