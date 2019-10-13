package web

import (
	"context"
	"fmt"
	"net/http"
)

func NewHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data, err := store.Fetch(ctx)
		if err != nil {
			return
		}
		fmt.Fprint(w, data)
	}
}

type Store interface {
	Fetch(ctx context.Context) (string, error)
}
