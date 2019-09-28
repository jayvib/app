package stack

import (
	"bufio"
	"io"
)

type CLI struct {
	stack Stack
	reader io.Reader
}

func (cli *CLI) Run() {
	scanner := bufio.NewScanner(cli.reader)
	scanner.Scan()
	cli.stack.Push(scanner.Text())
}
