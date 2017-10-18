//Package vozscrape, isolate package vozscrape from the rest of the application
package vozscrape

import (
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
func (g *Box) fetchThreads(Posts chan *Thread, pSelector *goquery.Selection, pLen int) {
	var wg sync.WaitGroup

	// We are telling to the WaitGroup
	// how many items should be in the channel
	// before it completes (the number is of course
	// the total of our Posts)
	wg.Add(pLen)

	// Loop through every Post in the page
	pSelector.Each(func(i int, s *goquery.Selection) {

		title := s.Find("h3")
		url, _ := title.Find("a").Attr("href")

		// Fetch every Post concurrently with
		// a GoRoutine
		go func(Threads chan *Thread) {
			defer wg.Done()
			Posts <- NewThreads(title.Text(), url)

		}(Posts)
	})

	wg.Wait()
}

func (g *Box) fetchBox() *scraper.Scraper {
	s := scraper.NewScraper(g.url, "utf-8")

	return s
}
