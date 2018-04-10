package main

import (
	// "fmt"

	"flag"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	"github.com/nnti3n/voz-archive-service/interface/requesthandler"
	// "github.com/pkg/errors"
)

var dev string

func init() {
	flag.StringVar(&dev, "dev", "true", "build for local dev")
}

func main() {
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

	log.Println("database_username", dbUser)
	log.Println("database_name", dbName)
	log.Println("database_password", dbPass)

	db := pg.Connect(&pg.Options{
		User:     dbUser,
		Database: dbName,
		Password: dbPass,
	})

	router := gin.Default()
	router.Use(cors.Default())
	env := requesthandler.Env{Db: db}
	r := router.Group("/")
	{
		// r.GET("/box", fetchAllBox)
		r.GET("box/:boxID", env.FetchAllThread)
		r.GET("thread/:threadID", env.FetchSingleThread)
		r.GET("thread/:threadID/posts", env.FetchThreadPosts)
	}
	router.Run()
}
