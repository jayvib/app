package main

import (
	"github.com/jayvib/app/cmd/hacking/testing/build-an-application/pocker"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error while opening %s %v", dbFileName, err)
	}

	store, err := pocker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem while initializing file system player store: %v",err)
	}
	svr := pocker.NewPlayerServer(store)
	if err := http.ListenAndServe(":5000", svr); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
