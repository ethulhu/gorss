package main

import (
	"fmt"
	"strings"
	"strconv"
	/* "io" */
	"os"
	_ "github.com/lib/pq"
	"database/sql"
	"encoding/json"
	"container/list"
	"github.com/Ezey/golor"
)

func count(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(*) FROM posts WHERE shown = false")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var unread string
		if err := rows.Scan(&unread); err != nil {
			panic(err)
		}
		fmt.Printf("%s unread posts\n",unread)
	}
}

type Post struct {
	Url string
	Title string
	Body string
}

func refresh(db *sql.DB) *list.List {
	rows, err := db.Query("SELECT title,url,body FROM posts WHERE unread = true")
	if err != nil {
		panic(err)
	}
	posts := list.New()
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title,&post.Url,&post.Body); err != nil {
			panic(err)
		}
		posts.PushBack(post)
	}
	return posts
}

func printPosts(posts *list.List) {
	i := 0
	for e := posts.Front(); e != nil; e = e.Next() {
		post := e.Value.(Post)

		fmt.Printf("%s: %s - %s\n",
				golor.Colorize(fmt.Sprintf("%d",i),golor.RED,golor.BLACK),
		    post.Title,
		    golor.Colorize(post.Url,golor.CYAN,golor.BLACK))
		i++
	}
}

func show(db *sql.DB) {
	posts := refresh(db)
	printPosts(posts)
	for true {
		var cmd string
		_, err := fmt.Scanln(&cmd)
		if err != nil { panic(err) }
		args := strings.Split(cmd," ")
		switch args[0] {
		case "q":
			os.Exit(0)
		case "r":
			posts := refresh(db)
			printPosts(posts)
		case "d","u","p":
			for i := range args[1:] {
				idx,_ := strconv.Atoi(args[i])
				e := posts.Front()
				for j := 0; j < idx; j++ {
					e := e.Next()
				}
				post := e.Value.(Post)
				fmt.Println(post.Body)
			}
		}
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage!!!")
		os.Exit(1)
	}

	conf_file, err := os.Open(os.ExpandEnv("$HOME/.rssrc"))
	if err != nil { panic(err) }
	buf := make([]byte, 1024)
	n, err := conf_file.Read(buf)
	if err != nil { panic(err) }
	var config interface{}
	err = json.Unmarshal(buf[:n], &config)
	if err != nil { panic(err) }
	var m map[string]interface{} = config.(map[string]interface{})
	m = m["db"].(map[string]interface{})
	var conninfo string
	conninfo = "user=" + m["user"].(string) + " dbname=" + m["database"].(string) + " host=" + m["host"].(string) + " password=" + m["password"].(string)

	if err := conf_file.Close(); err != nil { panic(err) }

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	switch os.Args[1] {
		case "count":
			count(db)
		case "show":
			show(db)
	}

}
