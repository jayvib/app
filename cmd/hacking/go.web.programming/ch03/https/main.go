package main

import (
	"fmt"
	"log"
	"net/http"
)

func main () {

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Hello SSL world!")
	}))

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}


	if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatal(err)
	}
}
