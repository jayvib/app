package pocker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T){
	store := NewInMemoryStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	t.Run("get score", func(t *testing.T){
		server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player, nil))
		server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player, nil))
		server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player, nil))

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newScoreGetRequest(player))

		assertScore(t, response, "3")
	})

	t.Run("get league", func(t *testing.T){
		response := httptest.NewRecorder()
		request := newLeagueRequest()
		server.ServeHTTP(response, request)

		want := []Player{
			{"Pepper", 3},
		}

		assertStatusCode(t, response, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertLeague(t, response, want)
	})
}