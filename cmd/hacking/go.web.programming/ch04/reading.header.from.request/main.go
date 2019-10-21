package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = `
<html>
<head>
<title>Template</title>
</head>
<body>
<ul>
	{{ range $key, $value := . }}
		<li><strong>{{ $key }}</strong> <strong>{{ $value }}</strong></li>
	{{ end }}
</ul>
</body>
</html>
`

func printEncoding(w http.ResponseWriter, r *http.Request) {
	encoding := r.Header.Get("Accept-Encoding")
	fmt.Fprintf(w, "Accept-Encoding: %v", encoding)
}

func main() {
	t := template.New("template")
	t = template.Must(t.Parse(tmpl))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		header := r.Header

		err := t.Execute(w, header)
		if err != nil {
			log.Println("Error:", err)
		}

	})

	http.HandleFunc("/encoding", printEncoding)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Hello there!")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
