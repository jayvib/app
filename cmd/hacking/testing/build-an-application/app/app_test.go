package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd": 10,
		},
	}
	svr := NewPlayerServer(store)
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
		svr := NewPlayerServer(store)
		response := httptest.NewRecorder()
		request := newWinPostRequest("Pepper", nil)
		svr.ServeHTTP(response, request)
		assertStatusCode(t, response, http.StatusAccepted)
	})

	t.Run("recording the win of Pepper", func(t *testing.T){
		store := &StubPlayerStore{
			scores: make(map[string]int),
		}
		svr := NewPlayerServer(store)
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

func TestLeague(t *testing.T) {
	t.Run("it returns 200 on /league", func(t *testing.T){
		store := &StubPlayerStore{
			scores: make(map[string]int),
		}
		server := NewPlayerServer(store)
		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatusCode(t, response, http.StatusOK)
	})

	t.Run("returns a league table as JSON", func(t *testing.T){
		league := []Player{
			{"Luffy", 10},
			{"Sanji", 11},
			{"Zoro", 12},
		}

		store := &StubPlayerStore{
			scores: make(map[string]int),
			league: league,
		}
		server := NewPlayerServer(store)
		request := newLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertLeague(t, response, league)
		assertContentType(t, response, jsonContentType)
	})
}

func newLeagueRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/league", nil)
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	contentType := response.Header().Get("Content-Type")
	if contentType != jsonContentType {
		t.Errorf("want '%s' got '%s'", "application/json", contentType)
	}
}

func assertLeague(t *testing.T, response *httptest.ResponseRecorder, want []Player) {
	t.Helper()
	var got []Player
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want '%v' got '%v'", want, got)
	}
}

func assertStatusCode(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
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
	league []Player
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
