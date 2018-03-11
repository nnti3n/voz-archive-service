//Package vozscrape isolate package vozscrape from the rest of the application
package vozscrape

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-pg/pg"
	"github.com/nnti3n/voz-archive-service/serviceWorker/scraper"
	"github.com/nnti3n/voz-archive-service/utilities"
)

var db = pg.Connect(&pg.Options{
	User:     "nntien",
	Database: "vozarchive",
})

// Box is the model for the forum box
type Box struct {
	ID      int
	Threads []*Thread `sql:"-"`
}

// NewBox Loop through every item with a goRoutine
func NewBox(boxPage int) *Box {

	b := new(Box)
	b.ID = 33
	b.Threads = b.getAsyncThreads(boxPage)

	return b
}

// Get all the Posts in the page with a goRoutine
func (b *Box) getAsyncThreads(boxPage int) []*Thread {

	// Start by scraping the forum box
	s := b.fetchBox(boxPage)

	// Slice of thread selector
	tPageSelector := []goquery.Selection{}
	// Count how many Thread there are in the page
	tLen := 0

	// Select all the threads in a box
	for _, page := range s {
		threadSelector := page.Find("#threadbits_forum_33 > tr")
		tLen += threadSelector.Size()
		tPageSelector = append(tPageSelector, *threadSelector)
	}

	fmt.Println("Number of threads ", tLen)

	// This is the slice that will contain all our Threads
	var t []*Thread

	// Construct a slice of Post chan, make it
	// of the size of the Posts found in the page
	tChan := make(chan *Thread, tLen)

	b.fetchThreads(tChan, tPageSelector, tLen)

	// Close the channel when the previous function is done
	close(tChan)

	for i := range tChan {
		t = append(t, i)
	}

	return t
}

// The function that will launch our GoRoutines: in order to prevent race conditions
// we will sync all of our routines
func (b *Box) fetchThreads(Threads chan *Thread, tPageSelector []goquery.Selection, pLen int) {
	var wg sync.WaitGroup

	// We are telling to the WaitGroup
	// how many items should be in the channel
	// before it completes (the number is of course
	// the total of our Threads)
	wg.Add(pLen)

	// Loop through every Thread in the page
	for _, tSelector := range tPageSelector {
		tSelector.Each(func(i int, s *goquery.Selection) {

			title := s.Find("td:nth-child(2) > div:first-child > a:last-of-type")
			// userID := -1
			_userID, exist := s.Find("td:nth-child(2) > .smallfont span:last-child").Attr("onclick")
			rUserID := regexp.MustCompile(`[\d]+`)
			userID, _ := strconv.Atoi(rUserID.FindString(_userID))

			userName := s.Find("td:nth-child(2) > .smallfont span").Text()

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
			postCount := s.Find("td:nth-child(4) a").Text()
			viewCount := s.Find("td:nth-child(5)").Text()

			// Fetch every Thread concurrently with
			// a GoRoutine
			go func(Threads chan *Thread) {
				defer wg.Done()
				Threads <- NewThread(utilities.ParseThreadURL(id), title.Text(), userID, userName, source, pageCount, postCount, viewCount, b.ID)

			}(Threads)
		})
	}

	wg.Wait()
}

func (b *Box) fetchBox(boxPage int) []scraper.Scraper {
	s := []scraper.Scraper{}
	for i := 1; i <= boxPage; i++ {
		t := scraper.NewScraper("https://vozforums.com/forumdisplay.php?f="+strconv.Itoa(b.ID)+"&page="+strconv.Itoa(i), "utf-8")
		// log.Println(t)
		s = append(s, *t)
	}

	return s
}
