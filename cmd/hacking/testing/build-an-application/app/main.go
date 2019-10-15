package main

import (
	"log"
	"net/http"
)

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{score: make(map[string]int)}
}

type InMemoryStore struct {
	score map[string]int
}

func (s *InMemoryStore) GetPlayerScore(player string) int {
	score, ok := s.score[player]
	if !ok {
		return 0
	}
	return score
}

func (s *InMemoryStore) RecordWin(player string) {
	s.score[player]++
}

func (s *InMemoryStore) GetLeague() []Player {
	var league []Player
	for name, win := range s.score {
		league = append(league, Player{Name: name, Wins: win})
	}
	return league
}

func main() {
	store := NewInMemoryStore()
	svr := &PlayerServer{store:store}
	if err := http.ListenAndServe(":5000", svr); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
