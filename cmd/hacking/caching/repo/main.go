package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
)

type Repository struct {
	Name string `json:"name"`
}

type Result struct {
	Repos    []Repository
	IsCached bool
}

type Client interface {
	GetRepositories() (Result, error)
	Username() (username string)
}

type ClientFunc func() (Result, error)

func (c ClientFunc) GetRepositories() (Result, error) {
	return c()
}

type ClientMiddleware func(client Client) Client

func NewCachedClient(client Client, cache *cache.Cache) *CachedClient {
	return &CachedClient{
		Client: client,
		cache:  cache,
	}
}

type CachedClient struct {
	Client
	cache *cache.Cache
}

func (r *CachedClient) GetRepositories() (res Result, err error) {
	var found bool
	res, found = r.get()
	if found {
		res.IsCached = true
		return
	}

	res, err = r.Client.GetRepositories()
	if err != nil {
		return Result{}, err
	}

	r.set(res)
	return res, nil
}

func (r *CachedClient) set(res Result) {
	r.cache.Set(r.Username(), res, cache.DefaultExpiration)
}

func (r *CachedClient) get() (res Result, found bool) {
	val, found := r.cache.Get(r.Username())
	if !found {
		return
	}
	return val.(Result), found
}

func CachedClientMiddleware(cache *cache.Cache) ClientMiddleware {
	return func(c Client) Client {
		return &CachedClient{
			Client: c,
			cache:  cache,
		}
	}
}

type LoggedClient struct {
	Client
	w io.Writer
}

func (l *LoggedClient) GetRepositories() (Result, error) {
	fmt.Fprintf(l.w, "Fetching %s repo", l.Username())
	return l.Client.GetRepositories()
}


func LoggedClientMiddleware(w io.Writer) ClientMiddleware {
	return func(c Client) Client {
		return &LoggedClient{
			Client: c,
			w:      w,
		}
	}
}
