package main

import (
	"github.com/jayvib/app/cmd/hacking/testing/build-an-application/pocker"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, closer, err := pocker.NewFileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer closer()
	svr := pocker.NewPlayerServer(store)
	if err := http.ListenAndServe(":5000", svr); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
