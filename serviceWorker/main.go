package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"github.com/nnti3n/voz-archive-service/serviceWorker/vozscrape"
)

var dev string

func init() {
	flag.StringVar(&dev, "dev", "true", "build for local dev")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// this is the console application
func main() {
	news := vozscrape.NewBox(33, 3)
	random := vozscrape.NewBox(17, 3)
	InsertDB(news)
	InsertDB(random)
}

// InsertDB will map model to db
func InsertDB(box *vozscrape.Box) {
	flag.Parse()

	var dbUser, dbPass, dbName string

	if dev == "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbUser = os.Getenv("VOZ_DATABASE_USER")
	dbName = os.Getenv("VOZ_DATABASE_NAME")
	dbPass = os.Getenv("VOZ_DATABASE_PASSWORD")

	db := pg.Connect(&pg.Options{
		User:     dbUser,
		Database: dbName,
		Password: dbPass,
	})

	_, err := db.Model(box).
		OnConflict("DO NOTHING").Insert()
	for _, thread := range box.Threads {
		_, err = db.Model(thread).
			OnConflict("(id) DO UPDATE").
			Set("page_count = ?page_count, post_count = ?post_count, view_count = ?view_count").
			Insert()
		for index, post := range thread.Posts {
			_, err = db.Model(post).
				OnConflict("DO NOTHING").Insert()
			if index+1 == len(thread.Posts) {
				_, err = db.Model(thread).
					Set("last_updated = ?", post.Time).
					Where("id = ?id").
					Update()
			}
		}
	}
	if err != nil {
		panic(err)
	}
	log.Println("Done!", box.ID)
}
