package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const jsonContentType = "application/json"

type Player struct {
	Name string
	Wins int
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	svr := &PlayerServer{
		store: store,
	}
	router := http.NewServeMux()
	router.HandleFunc("/league", svr.leagueHandler)
	router.HandleFunc("/players/", svr.playerHandler)
	svr.Handler = router
	return svr
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func (s *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.store.GetLeague())
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", jsonContentType)
}

func (s *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := getPlayer(r)
	switch r.Method {
	case http.MethodGet:
		s.showScore(w, player)
	case http.MethodPost:
		s.processWin(w, player)
	}
}

func (s *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := s.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprint(w, score)
}

func (s *PlayerServer) processWin(w http.ResponseWriter, player string) {
	s.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func getPlayer(r *http.Request) string {
	return r.URL.Path[len("/players/"):]
}
