package main

import (
	"log"
	"net/http"
)

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

func main() {
	store := &InMemoryStore{
		score: map[string]int{
			"Luffy": 20,
			"Sanji": 10,
		},
	}
	svr := &PlayerServer{store:store}
	if err := http.ListenAndServe(":5000", svr); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
