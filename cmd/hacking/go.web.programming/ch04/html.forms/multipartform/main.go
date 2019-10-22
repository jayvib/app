package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		fmt.Fprintln(w, r.PostForm)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
