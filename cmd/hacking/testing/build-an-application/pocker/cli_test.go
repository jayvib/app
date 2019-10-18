package pocker_test

import (
	"github.com/jayvib/app/cmd/hacking/testing/build-an-application/pocker"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("Sanji wins", func(t *testing.T){
		playerStore := pocker.NewStubPlayerStore()
		in := strings.NewReader("Sanji wins\n")
		cli := pocker.NewCLI(playerStore, in)

		cli.PlayPocker()
		want := "Sanji"
		pocker.AssertPlayeWin(t, playerStore, want)
	})

	t.Run("Luffy wins", func(t *testing.T){
		playerStore := pocker.NewStubPlayerStore()
		in := strings.NewReader("Luffy wins\n")
		cli := pocker.NewCLI(playerStore, in)

		cli.PlayPocker()
		want := "Luffy"
		pocker.AssertPlayeWin(t, playerStore, want)
	})
}

