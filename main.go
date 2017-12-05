package main

import (
	"fmt"
	"log"

	"github.com/nnti3n/voz-archive-plus/utilities"
	"github.com/nnti3n/voz-archive-plus/vozscrape"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// this is the console application
func main() {
	fmt.Println(Box())
}

// Box will scape forum boxes
func Box() string {
	s := vozscrape.NewBox()
	res2B, _ := utilities.JSONMarshal(s, true)
	return string(res2B)
}
