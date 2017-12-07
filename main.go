package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nnti3n/voz-archive-plus/utilities"
	"github.com/nnti3n/voz-archive-plus/vozscrape"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// this is the console application
func main() {
	// connStr := "user=nntien dbname=pqgotest sslmode=verify-full"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(Box())
	f, err := os.Create("json.txt")
	check(err)

	content := Box()
	n3, err := f.WriteString(content)
	fmt.Print("wrote %d bytes\n", n3)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Box will scape forum boxes
func Box() string {
	s := vozscrape.NewBox()
	res2B, _ := utilities.JSONMarshal(s, true)
	return string(res2B)
}
