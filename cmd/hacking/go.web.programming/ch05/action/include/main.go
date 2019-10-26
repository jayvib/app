package main

import (
	"html/template"
	"log"
	"net/http"
)

var t = template.Must(template.ParseGlob("*.html"))

func index(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "tmpl.html", "Hello World!")
}

func main() {
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
