package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		t, _ := template.ParseFiles("tmpl.html")
		t.Execute(w, "Hello World!")
	})
	http.ListenAndServe(":8080", nil)
}
