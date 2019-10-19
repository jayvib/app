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
	store, closer, err := pocker.NewFileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer closer()
	cli := pocker.NewCLI(store, os.Stdin, pocker.BlindAlerterFunc(pocker.StdOutAlerter))
	cli.PlayPocker()
}
