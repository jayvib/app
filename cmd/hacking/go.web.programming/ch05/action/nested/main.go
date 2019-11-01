package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var t = template.Must(template.ParseGlob("*.html"))

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Calling index")
	t.ExecuteTemplate(w, "layout",nil)
}

func main() {
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
