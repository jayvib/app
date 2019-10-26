package main

import (
	"html/template"
	"net/http"
)

var t = template.Must(template.ParseFiles("tmpl.html"))

func index(w http.ResponseWriter, r *http.Request) {
	l := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	t.Execute(w, l)
}

func fallback(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, nil)
}
func main() {
	http.HandleFunc("/iterate", index)
	http.HandleFunc("/fallback", fallback)
	http.ListenAndServe(":8080", nil)
}
