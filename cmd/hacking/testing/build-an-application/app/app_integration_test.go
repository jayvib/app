package main

import (
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T){
	store := NewInMemoryStore()
	server := PlayerServer{store: store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player, nil))
	server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player, nil))
	server.ServeHTTP(httptest.NewRecorder(), newWinPostRequest(player, nil))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newScoreGetRequest(player))

	assertScore(t, response, "3")
}