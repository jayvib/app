package main

import (
	"fmt"
	"log"
	"net/http"
)


func Log(h http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Printf("METHOD: %s PATH: %s\n", r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {

	http.Handle("/", Log(http.HandlerFunc(printnfoHandler)))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func printnfoHandler(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World, %s!", r.URL.Path)
}
