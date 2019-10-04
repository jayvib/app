package main

import (
	"html/template"
	"log"
	"net/http"
)

type Data struct {
	People []Person
}

type Person struct {
	Firstname string
	Lastname string
}

var people = []Person{
	{"Luffy", "Monkey"},
	{"Sanji", "Vinsmoke"},
	{"Zoro", "Roronoa"},
	{"Chopper", "Chopper"},
}

var tmpl = `
<html>
<header>
<title>Names</title>
</header>
<body>
<h1>Characters</h1>
<ul>
{{ range .People }}
	<li>{{ .Firstname }} {{ .Lastname }}</li>	
{{ end }}
</ul>
</body>
</html>`

func main() {
	tmp := template.Must(template.New("index").Parse(tmpl))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		err := tmp.Execute(w, &Data{People: people})
		if err != nil {
			log.Fatal(err)
		}
	})

	http.ListenAndServe(":8080", nil)
}
