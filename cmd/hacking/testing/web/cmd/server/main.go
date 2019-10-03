package main

import (
	"log"
	"net/http"
	"github.com/jayvib/app/cmd/hacking/testing/web"
)

func main() {
	handler := http.HandlerFunc(web.PlayerServer)
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
