package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/nnti3n/voz-archive-plus/serviceWorker/vozscrape"
	"github.com/nnti3n/voz-archive-plus/utilities"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// this is the console application
func main() {
	s := vozscrape.NewBox(5)
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

	_, err := db.Model(box).
		OnConflict("DO NOTHING").Insert()
	for _, thread := range box.Threads {
		_, err = db.Model(thread).
			OnConflict("DO NOTHING").Insert()
		err = db.Insert(thread)
		for _, post := range thread.Posts {
			_, err = db.Model(post).
				OnConflict("DO NOTHING").Insert()
		}
	}
	if err != nil {
		panic(err)
	}
}
