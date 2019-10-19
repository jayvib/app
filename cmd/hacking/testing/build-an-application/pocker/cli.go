package pocker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const PlayerPrompt= "Please enter the number of players: "

type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)
func(b BlindAlerterFunc) ScheduledAlertAt(duration time.Duration, amount int) {
	b(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func(){
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}

func NewCLI(store PlayerStore, input io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		store: store,
		scanner: bufio.NewScanner(input),
		blindAlerter: alerter,
		out: out,
	}
}

type CLI struct {
	store PlayerStore
	scanner *bufio.Scanner
	blindAlerter BlindAlerter
	out io.Writer
}

func (c *CLI) PlayPocker() {
	fmt.Fprint(c.out, PlayerPrompt)
	c.scheduledBlindAlerts()
	c.store.RecordWin(extractName(c.readline()))
}

func (c *CLI) scheduledBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		c.blindAlerter.ScheduledAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}

func (c *CLI) readline() string {
	c.scanner.Scan()
	return c.scanner.Text()
}

func extractName(s string) string {
	return strings.Replace(s, " wins", "", 1)
}