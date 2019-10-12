package counter

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T){
		counter := &Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()
		assertCount(t, counter, 3)
	})

	t.Run("it runs concurrently safe", func(t *testing.T){
		workers := 1000
		var wg sync.WaitGroup

		counter := &Counter{}

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup)	{
				defer wg.Done()
				counter.Inc()
			}(&wg)
		}

		wg.Wait()

		assertCount(t, counter, workers)
	})
}

func assertCount(t *testing.T, counter *Counter, want int) {
	t.Helper()
	got := counter.Value()
	if got != want {
		t.Errorf("got '%d' want '%d'", got, 3)
	}
}