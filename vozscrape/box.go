//Package vozscrape isolate package vozscrape from the rest of the application
package vozscrape

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/nnti3n/voz-archive-plus/scraper"
	"github.com/nnti3n/voz-archive-plus/utilities"
)

// Box is the model for the forum box
type Box struct {
	url string
	id  int

	Threads []*Thread
}

// NewBox Loop through every item with a goRoutine
func NewBox() *Box {

	g := new(Box)
	g.url = "https://vozforums.com/forumdisplay.php?f=33"
	fmt.Println(g.url)
	g.Threads = g.getAsyncThreads()
	log.Println("G ", g)

	return g
}

// Get all the Posts in the page with a goRoutine
func (g *Box) getAsyncThreads() []*Thread {

	// Start by scraping the forum box
	s := g.fetchBox()

	// Select all the threads in a box
	tSelector := s.Find("#threadbits_forum_33 > tr")

	// Count how many Thread there are in the page
	tLen := tSelector.Size()
	fmt.Println("Number of threads ", tLen)

	// This is the slice that will contain all our Threads
	var p []*Thread

	// Construct a slice of Post chan, make it
	// of the size of the Posts found in the page
	tChan := make(chan *Thread, tLen)

	g.fetchThreads(tChan, tSelector, tLen)

	// Close the channel when the previous function is done
	close(tChan)

	for i := range tChan {
		p = append(p, i)
	}

	return p
}

// The function that will launch our GoRoutines: in order to prevent race conditions
// we will sync all of our routines
func (g *Box) fetchThreads(Threads chan *Thread, pSelector *goquery.Selection, pLen int) {
	var wg sync.WaitGroup

	// We are telling to the WaitGroup
	// how many items should be in the channel
	// before it completes (the number is of course
	// the total of our Posts)
	wg.Add(pLen)

	// Loop through every Post in the page
	pSelector.Each(func(i int, s *goquery.Selection) {

		title := s.Find("td:nth-child(2) > div:first-child > a:last-of-type")
		source, exist := s.Find("td:nth-child(2)").Attr("title")
		if !exist {
			log.Println("Not found source ", exist)
		}
		id, exist := title.Attr("href")
		if !exist {
			log.Println("Not found id", exist)
		}
		pageCount := "1"
		pageURL, exist := s.Find("td:nth-child(2) div:first-child span.smallfont a:last-child").Attr("href")
		if !exist {
			log.Println("Not found pageURL PageCount 1")
		} else {
			pageCount = strings.Split(pageURL, "page=")[1]
		}

		// fmt.Print("ID ", utilities.ParseThreadURL(id), " ")
		// fmt.Print("title ", title.Text(), " ")
		// fmt.Print("source ", source, " ")
		// fmt.Print("pageCount ", pageCount, "\n", "\n")

		// Fetch every Thread concurrently with
		// a GoRoutine
		go func(Threads chan *Thread) {
			defer wg.Done()
			Threads <- NewThread(utilities.ParseThreadURL(id), title.Text(), source, pageCount)

		}(Threads)
	})

	wg.Wait()
}

func (g *Box) fetchBox() *scraper.Scraper {
	s := scraper.NewScraper(g.url, "utf-8")
	fmt.Println(s)

	return s
}
