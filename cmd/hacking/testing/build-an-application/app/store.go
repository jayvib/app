package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type PlayerStore interface {
	GetPlayerScore(player string) int
	RecordWin(player string)
	GetLeague() League
}

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

func (s *InMemoryStore) GetLeague() League {
	var league []Player
	for name, win := range s.score {
		league = append(league, Player{Name: name, Wins: win})
	}
	return league
}

func NewFileSystemPlayerStore(database *os.File) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{database}),
		league:   league,
	}
}

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func (f *FileSystemPlayerStore) GetLeague() (league League) {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.GetLeague().Find(name)
	if player == nil {
		return 0
	}
	return player.Wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{Name: name, Wins: 1})
	}

	err := f.database.Encode(f.league)
	if err != nil {
		log.Println(err)
	}
}

func NewLeague(r io.Reader) (League, error) {
	var league []Player
	err := json.NewDecoder(r).Decode(&league)
	if err != nil {
		return nil, fmt.Errorf("problem while parsing league: %v", err)
	}
	return league, nil
}
