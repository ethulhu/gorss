package main

import (
	"bufio"
	"os"
	"os/exec"
	"fmt"
	"strings"
	"strconv"
	"database/sql"
	"github.com/Ezey/golor"
)

type Post struct {
	Url string
	Title string
	Body string
}

func refresh(db *sql.DB) []Post {
	rows, err := db.Query("SELECT title,url,body FROM posts WHERE unread = true")
	if err != nil {
		panic(err)
	}
	posts := make([]Post,0,10)
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title,&post.Url,&post.Body); err != nil {
			panic(err)
		}
		posts = append(posts,post)
	}
	return posts
}

func printPosts(posts []Post) {
	for i := 0; i < len(posts); i++ {
		post := posts[i]

		fmt.Printf("%s: %s - %s\n",
				golor.Colorize(fmt.Sprintf("%d",i),golor.RED,golor.BLACK),
		    post.Title,
		    golor.Colorize(post.Url,golor.CYAN,golor.BLACK))
	}
}

func show(db *sql.DB) {
	posts := refresh(db)
	printPosts(posts)
	fmt.Print(":")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		args := strings.Split(cmd," ")
		switch args[0] {
		case "q":
			return
		case "r":
			posts := refresh(db)
			printPosts(posts)
		case "d","u","p":
			for _,r := range args[1:] {
				rang := argParse(r)
				for _,i := range rang {
					switch args[0] {
					case "p":
						cmd := exec.Command("less")
						cmd.Stdin = strings.NewReader(reverseMarkdown(posts[i].Body))
						cmd.Stdout = os.Stdout
						cmd.Run()
					case "d","u":
						fmt.Println("todo")
					}
				}
			}
		default:
			fmt.Println("Usage: (q|r|d|u|p) <range>")
		}
		fmt.Print(":")
	}
}

func argParse(arg string) []int {
	args := strings.Split(arg,"-")
	switch len(args) {
	case 1:
		num, err := strconv.Atoi(args[0])
		if err != nil { return nil }
		nums := make([]int,1)
		nums[0] = num
		return nums
	case 2:
		num1, err := strconv.Atoi(args[0])
		if err != nil { return nil }
		num2, err := strconv.Atoi(args[1])
		if err != nil { return nil }
		if num1 > num2 {
			num1, num2 = num2, num1
		}
		nums := make([]int,(num2-num1)+1)
		for i := 0; i < (num2-num1)+1; i++ {
			nums[i] = num1 + i
		}
		return nums
	default:
		return nil
	}
	return nil
}

func reverseMarkdown(html string) string {
	cmd := exec.Command("reverse_markdown")
	cmd.Stdin = strings.NewReader(html)
	markdown, err := cmd.Output()
	if err != nil { panic(err) }
	return string(markdown)
}
