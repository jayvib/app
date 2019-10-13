package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	ctx context.Context
	response string
	cancelCalled bool
	t *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	s.ctx = ctx
	data := make(chan string, 1)
	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				s.t.Log("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}

func (s *SpyStore) Cancel() {
	s.cancelCalled = true
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelCalled {
		s.t.Errorf("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelCalled {
		s.t.Errorf("should not be cancelled, but was cancelled")
	}
}

type SpyResponseWriter struct {
	written bool
}

func (r *SpyResponseWriter) Header() http.Header {
	r.written = true
	return nil
}

func (r *SpyResponseWriter) Write([]byte) (int, error) {
	r.written = true
	return 0, errors.New("not implemented")
}

func (r *SpyResponseWriter) WriteHeader(statusCode int) {
	r.written = true
}

func TestHandler(t *testing.T) {
	data := "hello, world"
	t.Run("server returns an 'hello, world' response", func(t *testing.T){
		stubStore := &SpyStore{response: data, t: t}

		svr := NewHandler(stubStore)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf("got '%s', want '%s'", response.Body.String(), data)
		}
		stubStore.assertWasNotCancelled()
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T){
		store := &SpyStore{response: data, t: t}

		svr := NewHandler(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(request.Context())

		time.AfterFunc(5 * time.Millisecond, cancel)

		// I need to make sure that when an error occur where should no writes
		// to the response writer
		request = request.WithContext(cancellingCtx)

		response := &SpyResponseWriter{}
		svr.ServeHTTP(response, request)

		if store.ctx != request.Context() {
			t.Errorf("context was not passed to store")
		}

		if response.written {
			t.Error("a response should not have been written")
		}
	})
}
