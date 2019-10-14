package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd": 10,
		},
	}
	svr := &PlayerServer{store: store}
	t.Run("returns Pepper's score", func(t *testing.T){
		request := newScoreGetRequest("Pepper")
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)
		assertScore(t, response, "20")
		assertStatusCode(t, response, http.StatusOK)
	})

	t.Run("returns Floyd's score", func(t *testing.T){
		request := newScoreGetRequest("Floyd")
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)
		assertScore(t, response, "10")
		assertStatusCode(t, response, http.StatusOK)
	})

	t.Run("returns 404 on missing players", func(t *testing.T){
		request := newScoreGetRequest("Luffy")
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)
		assertStatusCode(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T){
	t.Run("accepts POST request", func(t *testing.T){
		store := &StubPlayerStore{
			scores: make(map[string]int),
		}
		svr := &PlayerServer{store:store}
		response := httptest.NewRecorder()
		request := newWinPostRequest("Pepper", nil)
		svr.ServeHTTP(response, request)
		assertStatusCode(t, response, http.StatusAccepted)
	})

	t.Run("recording the win of Pepper", func(t *testing.T){
		store := &StubPlayerStore{
			scores: make(map[string]int),
		}
		svr := &PlayerServer{store:store}
		response := httptest.NewRecorder()
		request := newWinPostRequest("Pepper", nil)
		svr.ServeHTTP(response, request)
		assertStatusCode(t, response, http.StatusAccepted)
		if len(store.winCalls) != 1 {
			t.Errorf("expecting call to be %d but got %d", 1, len(store.winCalls))
		}
		if store.winCalls[0] != "Pepper"  {
			t.Errorf("expecting Pepper but found %s", store.winCalls[0])
		}
	})
}

func assertStatusCode(t *testing.T, response *httptest.ResponseRecorder, want int) {
	if response.Code != want {
		t.Errorf("expecting status code '%d' but got '%d'", want, response.Code)
	}
}

func assertScore(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Body.String()
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func newWinPostRequest(player string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost,fmt.Sprintf( "/players/%s", player), body)
	return request
}

func newScoreGetRequest(player string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return req
}

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
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
