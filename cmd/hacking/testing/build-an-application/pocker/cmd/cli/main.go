package main

import (
	"fmt"
	"github.com/jayvib/app/cmd/hacking/testing/build-an-application/pocker"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play pocker")
	fmt.Println("Type '{Name} wins' to record a win")
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("unable to open file:", err)
	}
	defer db.Close()

	store, err := pocker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatal("problem while initializing file system play store:", err)
	}

	cli := pocker.NewCLI(store, os.Stdin)
	cli.PlayPocker()
}
