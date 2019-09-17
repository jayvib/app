package elasticsearch

import (
	"context"
	"fmt"
	"github.com/jayvib/app/config"
	"github.com/olivere/elastic/v7"
	"sync"
)

var (
	once   sync.Once
	client *elastic.Client
)

var (
	onceSimpleClient sync.Once
	simpleClient     *elastic.Client
)

func NewClient() (c *elastic.Client, err error) {
	var conf *config.Config
	once.Do(func() {
		conf, err = config.New()
		if err != nil {
			return
		}
		client, err = elastic.NewClient(
			elastic.SetURL(conf.Elasticsearch.Servers...))
		if err != nil {
			return
		}
		ping(client, conf.Elasticsearch.Servers...)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewSimpleClient() (c *elastic.Client, err error) {
	var conf *config.Config
	onceSimpleClient.Do(func() {
		conf, err = config.New()
		if err != nil {
			return
		}
		simpleClient, err = elastic.NewSimpleClient(
			elastic.SetURL(conf.Elasticsearch.Servers...))
		if err != nil {
			return
		}
		ping(simpleClient, conf.Elasticsearch.Servers...)
	})
	if err != nil {
		return nil, err
	}
	return simpleClient, nil
}

func ping(c *elastic.Client, urls ...string) (err error) {
	fmt.Println("==========ELASTICSEARCH=========")
	for _, url := range urls {
		fmt.Printf("%s: PING...", url)
		_, _, err = c.Ping(url).Do(context.Background())
		if err != nil {
			fmt.Print("NOT OK")
		} else {
			fmt.Print("OK")
		}
		fmt.Print("\n")
	}
	if err != nil {
		return
	}
	return nil
}
