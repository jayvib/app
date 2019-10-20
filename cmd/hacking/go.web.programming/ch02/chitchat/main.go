package main

import (
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("assets"))
	mux.Handle("/static/", files)
	mux.Handle("/", http.HandlerFunc(index))

	server := &http.Server{
		Addr: "0.0.0.0:8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
