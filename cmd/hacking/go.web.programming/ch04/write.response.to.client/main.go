package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func display(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<html>
<head>
<title>ResponseWriter</title>
</head>
<body>
<p>This is an example of writing an HTML content to the response writer</p>
</body>
</html>`
	fmt.Fprint(w, tmpl)
}

func writeHeader(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "No such serice, try next door")
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(http.StatusFound)
}

func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	post := struct {
		User string `json:"user"`
		Threads []string `json:"threads"`
	}{
		User: "Luffy Monkey",
		Threads: []string{"first", "second", "third"},
	}
	json.NewEncoder(w).Encode(post)
}

func main() {
	http.HandleFunc("/", display)
	http.HandleFunc("/writeHeader", writeHeader)
	http.HandleFunc("/header", headerHandler)
	http.HandleFunc("/json", jsonExample)
	http.ListenAndServe(":8080", nil)
}
