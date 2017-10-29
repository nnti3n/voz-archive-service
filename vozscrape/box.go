//Package vozscrape, isolate package vozscrape from the rest of the application
package vozscrape

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/ref/mercato/scraper"
)

// Box is the model for the forum box
type Box struct {
	url string
	id  int

	Threads []*Thread `json:"results"`
}

// NewBox Loop through every item with a goRoutine
func NewBox() *Box {

	g := new(Box)
	g.url = "https://vozforums.com/forumdisplay.php?f=33"
	g.Threads = g.getAsyncThreads()

	return g
}

// Get all the Posts in the page with a goRoutine
func (g *Box) getAsyncThreads() []*Thread {

	// Start by scraping the forum box
	s := g.fetchBox()

	// Select all the threads in a box
	pSelector := s.Find("#33 > tr")

	// Count how many Thread there are in the page
	pLen := pSelector.Size()

	// This is the slice that will contain all our Threads
	var p []*Thread

	// Construct a slice of Post chan, make it
	// of the size of the Posts found in the page
	pChan := make(chan *Thread, pLen)

	g.fetchThreads(pChan, pSelector, pLen)

	// Close the channel when the previous function is done
	close(pChan)

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
		source, _ := s.Find("td:nth-child(2)").Attr("title")
		fmt.Println(title)
		id, _ := title.Attr("href")
		pageURL, _ := s.Find("td:nth-child(2) div:first-child span.smallfont a:last-child").Attr("href")
		pageCount := "1"
		if pageURL != "" {
			pageCount = strings.Split(pageURL, "page=")[1]
		}

		// Fetch every Thread concurrently with
		// a GoRoutine
		go func(Threads chan *Thread) {
			defer wg.Done()
			Threads <- NewThread(id, title.Text(), source, pageCount)

		}(Threads)
	})

	wg.Wait()
}

func (g *Box) fetchBox() *scraper.Scraper {
	s := scraper.NewScraper(g.url, "utf-8")

	return s
}
