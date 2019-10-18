package main

import (
	"bytes"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Requirement:
// - Cache the repository result
// when it first fetch.
// - When the result is already cached
// then return the cached result.

type stubClient struct {
	username string
	result   Result
}

func (s stubClient) GetRepositories() (Result, error) {
	return s.result, nil
}

func (s stubClient) Username() string {
	return s.username
}

func assertResult(t *testing.T, got Result, want Result) {
	assert.Equal(t, want, got)
}


func TestCachedGetRepositories(t *testing.T) {
	cache := cache.New(1*time.Minute, 1*time.Minute)

	t.Run("return cached result after calling another GetRepositories", func(t *testing.T) {
		dummyRepos := []Repository{
			{"github.com/jayvib/pocker"},
			{"github.com/luffy/pocker"},
			{"github.com/sanji/pocker"},
		}

		stub := stubClient{
			username: "testuser",
			result: Result{
				Repos: dummyRepos,
			},
		}

		cachedGetRepo := NewCachedClient(stub, cache)

		t.Run("first call to GetRepositories should not from cache", func(t *testing.T){
			got, _ := cachedGetRepo.GetRepositories()
			want := Result{
				Repos: dummyRepos,
			}
			assertResult(t, got, want)
		})

		t.Run("succeding call to cached client should be from cache", func(t *testing.T){
			got, _ := cachedGetRepo.GetRepositories()
			want := Result{
				Repos: dummyRepos,
				IsCached: true,
			}
			assertResult(t, want, got)
		})

	})

	t.Run("return cached result from another client", func(t *testing.T) {
		dummyRepos := []Repository{
			{"github.com/nami/cloud"},
			{"github.com/nami/taktiile"},
			{"github.com/nami/thunderbold"},
		}

		stub := stubClient{
			username: "nami",
			result: Result{
				Repos: dummyRepos,
			},
		}

		cachedGetRepo := NewCachedClient(stub, cache)

		t.Run("first call result should not be from cache", func(t *testing.T) {
			got, _ := cachedGetRepo.GetRepositories()
			want := Result{
				Repos: dummyRepos,
			}
			assertResult(t, want, got)
		})

		t.Run("second call result should be from cache", func(t *testing.T) {
			got, _ := cachedGetRepo.GetRepositories()
			want := Result{
				Repos: dummyRepos,
				IsCached: true,
			}
			assertResult(t, want, got)
		})
	})

	t.Run("decorating the client with caching middleware", func(t *testing.T){
		dummyRepos := []Repository{
			{"github.com/sanji/cloud"},
			{"github.com/sanji/taktiile"},
			{"github.com/sanji/thunderbold"},
		}

		stub := stubClient{ // this is my main client
			username: "sanji",
			result: Result{
				Repos: dummyRepos,
			},
		}

		cachedGetRepo :=  CachedClientMiddleware(cache)(stub)

		t.Run("first call result should not be from cache", func(t *testing.T) {
			got, _ := cachedGetRepo.GetRepositories()
			want := Result{
				Repos: dummyRepos,
			}
			assertResult(t, want, got)
		})

		t.Run("second call result should be from cache", func(t *testing.T) {
			got, _ := cachedGetRepo.GetRepositories()
			want := Result{
				Repos: dummyRepos,
				IsCached: true,
			}
			assertResult(t, want, got)
		})

	})
}

func TestLoggedClient(t *testing.T) {
	var buff bytes.Buffer

	dummyRepos := []Repository{
		{"github.com/sanji/cloud"},
		{"github.com/sanji/taktiile"},
		{"github.com/sanji/thunderbold"},
	}

	stub := stubClient{ // this is my main client
		username: "sanji",
		result: Result{
			Repos: dummyRepos,
		},
	}
	loggedClient := LoggedClientMiddleware(&buff)(stub)
	got, _ := loggedClient.GetRepositories()
	want := Result{
		Repos: dummyRepos,
	}
	assert.Equal(t, want, got)
	assert.Equal(t, "Fetching sanji repo", buff.String())
}
