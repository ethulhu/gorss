package main

import (
	"fmt"
	/* "io" */
	"os"
	_ "github.com/lib/pq"
	"database/sql"
	"encoding/json"
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

func show(db *sql.DB) {
	rows, err := db.Query("SELECT title,url,body FROM posts WHERE unread = true")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var title string
		var url string
		var body string
		if err := rows.Scan(&title,&url,&body); err != nil {
			panic(err)
		}
		fmt.Printf("%s ==> %s\n",title,url)
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
