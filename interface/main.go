package main

import (
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/nnti3n/voz-archive-plus/interface/requesthandler"
	// "github.com/pkg/errors"
)

// func LoadConfiguration(pwd string) error {
// 	viper.SetConfigName("voz-config")
// 	viper.AddConfigPath(pwd)
// 	devPath := pwd[:len(pwd)-3] + "interface/"
// 	_, file, _, _ := runtime.Caller(1)
// 	configPath := path.Dir(file)
// 	viper.AddConfigPath(devPath)
// 	viper.AddConfigPath(configPath)
// 	return viper.ReadInConfig() // Find and read the config file
// }

func main() {
	// flag.Parse()
	// pwd, err := osext.ExecutableFolder()
	// if err != nil {
	// 	log.Fatalf("cannot retrieve present working directory: %i", 0600, nil)
	// }

	// err = LoadConfiguration(pwd)
	// if err != nil {
	// 	panic(errors.Errorf("Fatal reading config file: %s \n", err))
	// }

	// var db *pg.DB
	// if dev == "true" {
	// 	db = pg.Connect(fmt.Sprintf("host=%s port=%d user=%s "+
	// 		"password=%s dbname=%s sslmode=disable",
	// 		dbURL, dbPort, dbUser, dbPass, dbName))
	// 	if err != nil {
	// 		panic(errors.Errorf("Cannot connect to database: %s", err))
	// 	}
	// } else {
	// 	db = pg.Connect(dbURL)
	// 	if err != nil {
	// 		panic(errors.Errorf("Cannot connect to database: %s", err))
	// 	}
	// }

	db := pg.Connect(&pg.Options{
		User:     "nntien",
		Database: "vozarchive",
	})

	router := gin.Default()
	env := requesthandler.Env{Db: db}
	r := router.Group("/api")
	{
		// r.GET("/box", fetchAllBox)
		r.GET("/box/:boxID", env.FetchAllThread)
		r.GET("/thread/:threadID", env.FetchSingleThread)
		r.GET("/thread/:threadID/posts", env.FetchThreadPosts)
	}
	router.Run()
}
