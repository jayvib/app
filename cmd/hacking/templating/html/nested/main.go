package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	tmpl := template.Must(template.ParseGlob("*.tmpl"))
	err := tmpl.ExecuteTemplate(os.Stdout, "index", nil)
	if err != nil {
		log.Fatal(err)
	}
}
