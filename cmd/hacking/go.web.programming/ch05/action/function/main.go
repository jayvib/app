package main

// This demonstrates how to use the pipeline in template.

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

var tmpl = `
<html>
<head>
<title>Template Functions</title>
</head>
<body>
<div>The date/time is {{ . | fdate }}</div>
</body>
</html>`

var tmpl2 = `
<html>
<head>
<title>Template Functions</title>
</head>
<body>
<div>Normal case {{ . }}</div>
<div>Upper case {{ toUpper . }}</div>
</body>
</html>`

func formatDate(t time.Time) string {
	layout := "2006-01-02"
	return t.Format(layout)
}

func handlerDateNow(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"fdate": formatDate,
	}
	t := template.Must(
		template.New("tmpl").
			Funcs(funcMap).
			Parse(tmpl))
	t.Execute(w, time.Now())
}

func handleUpper(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}
	t := template.Must(
		template.New("tmpl").
			Funcs(funcMap).
			Parse(tmpl2))
	t.Execute(w, "luffy")
}

func main() {
	http.HandleFunc("/date", handlerDateNow)
	http.HandleFunc("/upper", handleUpper)
	http.ListenAndServe(":8080", nil)
}
