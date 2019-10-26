package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

var t = template.Must(template.ParseFiles("tmpl.html"))

func index(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	t.Execute(w, "Hey")
}
func fallback(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/fb", fallback)
	http.ListenAndServe(":8080", nil)
}
