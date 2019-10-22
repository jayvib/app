package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		fmt.Fprint(w, r.PostForm)
		fmt.Fprint(w, r.Form)
	})
	http.HandleFunc("/multipart", func(w http.ResponseWriter, r *http.Request){
		r.ParseMultipartForm(1024)
		fmt.Fprintln(w, r.MultipartForm)
	})
	http.ListenAndServe(":8080", nil)
}
