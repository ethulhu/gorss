package main

import (
	"database/sql"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
)

func add(db *sql.DB, urls []string) {
	for _, url := range urls {
		add_feed(db, url)
	}
}

func add_feed(db *sql.DB, url string) {
	chanHandler := func(feed *rss.Feed, channels []*rss.Channel) {
		_, err := db.Exec("INSERT INTO feeds (name, url) VALUES ($1, $2)", feed.Url, channels[0].Title)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s added\n", channels[0].Title)
	}
	feed := rss.New(10, true, chanHandler, nil)
	if err := feed.Fetch(url, nil); err != nil {
		panic(err)
	}
}
