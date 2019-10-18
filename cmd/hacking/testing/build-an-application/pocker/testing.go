package pocker

import "testing"

func NewStubPlayerStore() *StubPlayerStore {
	return &StubPlayerStore{scores: make(map[string]int)}
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	score, ok := s.scores[player]
	if !ok {
		return 0
	}
	return score
}

func (s *StubPlayerStore) RecordWin(player string) {
	s.winCalls = append(s.winCalls, player)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}
func AssertPlayeWin(t *testing.T, playerStore *StubPlayerStore, want string) {
	if len(playerStore.winCalls) != 1 {
		t.Fatalf("expecting win calls to be '%d' but got '%d'", 1, len(playerStore.winCalls))
	}

	if got := playerStore.winCalls[0]; got != want {
		t.Errorf("want '%s' got '%s'", want, got)
	}
}
