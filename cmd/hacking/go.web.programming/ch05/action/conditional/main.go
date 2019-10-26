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
	t.Execute(w, rand.Intn(10) > 5)
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
