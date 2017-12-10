package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/nnti3n/voz-archive-plus/utilities"
	"github.com/nnti3n/voz-archive-plus/vozscrape"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// this is the console application
func main() {
	s := vozscrape.NewBox()
	f, err := os.Create("json.txt")
	check(err)

	content := serialStruct(s)
	_, err = f.WriteString(content)
	check(err)
	fmt.Println("Wrote json.txt")

	DbModel(s)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Box will scape forum boxes
func serialStruct(s *vozscrape.Box) string {
	res2B, _ := utilities.JSONMarshal(s, true)
	return string(res2B)
}

// DbModel will map model to db
func DbModel(box *vozscrape.Box) {
	db := pg.Connect(&pg.Options{
		User:     "nntien",
		Database: "vozarchive",
	})

	err := db.Insert(box)
	for _, thread := range box.Threads {
		err = db.Insert(thread)
		for _, post := range thread.Posts {
			if post.ID == 0 {
				fmt.Println(post.Content)
			}
			err = db.Insert(post)
		}
	}
	if err != nil {
		panic(err)
	}
}
