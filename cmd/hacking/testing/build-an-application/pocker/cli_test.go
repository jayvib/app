package pocker_test

import (
	"fmt"
	"github.com/jayvib/app/cmd/hacking/testing/build-an-application/pocker"
	"strings"
	"testing"
	"time"
)

var dummyAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {
	t.Run("Sanji wins", func(t *testing.T){
		playerStore := pocker.NewStubPlayerStore()
		in := strings.NewReader("Sanji wins\n")
		cli := pocker.NewCLI(playerStore, in, dummyAlerter)

		cli.PlayPocker()
		want := "Sanji"
		pocker.AssertPlayeWin(t, playerStore, want)
	})

	t.Run("Luffy wins", func(t *testing.T){
		playerStore := pocker.NewStubPlayerStore()
		in := strings.NewReader("Luffy wins\n")
		cli := pocker.NewCLI(playerStore, in,dummyAlerter)

		cli.PlayPocker()
		want := "Luffy"
		pocker.AssertPlayeWin(t, playerStore, want)
	})

	t.Run("it schedules printing of blind values", func(t *testing.T){
		in := strings.NewReader("Chris wins\n")
		store := pocker.NewStubPlayerStore()
		alerter := &SpyBlindAlerter{}

		cli := pocker.NewCLI(store, in, alerter)
		cli.PlayPocker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(c.String(), func(t *testing.T){
				if len(alerter.alerts) <= i {
					t.Fatalf("alert %d is not scheduled %v", i, alerter.alerts)
				}

				// check for the amount
				alert := alerter.alerts[i]
				assertScheduledAlert(t, alert, c)
			})
		}
	})
}

func assertScheduledAlert(t *testing.T, got scheduledAlert, want scheduledAlert) {
	gotAmount := got.amount
	if gotAmount != want.amount {
		t.Errorf("amount want '%d' but got '%d'", want.amount, gotAmount)
	}
	gotScheduledAt := got.at
	if gotScheduledAt != want.at {
		t.Errorf("schedule want '%v', schedule got '%v'", want.at, gotScheduledAt)
	}
}

type scheduledAlert struct {
	at time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduledAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts,scheduledAlert{at: duration, amount: amount})
}
