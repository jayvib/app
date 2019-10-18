package pocker

import (
	"bufio"
	"io"
	"strings"
)

func NewCLI(store PlayerStore, input io.Reader) *CLI {
	return &CLI{
		store: store,
		scanner: bufio.NewScanner(input),
	}
}

type CLI struct {
	store PlayerStore
	scanner *bufio.Scanner
}

func (c *CLI) PlayPocker() {
	c.store.RecordWin(extractName(c.readline()))
}

func (c *CLI) readline() string {
	c.scanner.Scan()
	return c.scanner.Text()
}

func extractName(s string) string {
	return strings.Replace(s, " wins", "", 1)
}