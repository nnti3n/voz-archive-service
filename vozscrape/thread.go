package vozscrape

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ref/mercato/scraper"
)

// Thread is the model for the forums threads
type Thread struct {
	ID          int
	Title       string
	Source      string
	PageCount   int
	PostCount   int
	LastUpdated string

	Posts []*Post
}

// Post is the struct for a single Post in the Grocery Store
type Post struct {
	url      string
	PostID   int
	Number   int
	UserID   int
	UserName string
	Time     string
	Content  string
}

// NewThread creates a Thread and fills missing information
// from the Thread page
func NewThread(id int, title string, source string, pageCount string) *Thread {

	t := new(Thread)
	t.ID = id
	t.Title = title
	t.Source = source
	t.PageCount, _ = strconv.Atoi(pageCount)

	// Start scraping thread
	tPage := t.fetchThread()

	t.Posts = t.getPosts(tPage)

	return t
}

func (t *Thread) getPosts(pPage *scraper.Scraper) []*Post {
	posts := []*Post{}
	fmt.Println(pPage)
	fmt.Println("postLen", pPage.Find("#posts > div [align='left']").Size())

	pPage.Find("#posts > div [align='left']").Each(func(i int, s *goquery.Selection) {
		p := new(Post)
		number := s.Find("tr:first-child td div:first-child a:first-child")
		numberName, exist := number.Attr("name")
		if !exist {
			fmt.Println("Not found name")
		}
		p.Number, _ = strconv.Atoi(numberName)
		_number, _ := number.Attr("href")
		p.PostID, _ = strconv.Atoi(strings.Split(strings.Split(_number, "=")[1], "&")[0])

		p.Time = strings.TrimSpace(s.Find("tr:first-child td.thead div:nth-child(2)").Text())

		username, _ := s.Find(".bigusername").Attr("href")
		p.UserName = strings.Split(username, "u=")[1]

		p.Content = strings.TrimSpace(s.Find(".voz-post-message").Text())

		// fmt.Print("p.PostID ", p.PostID)
		// fmt.Print(" p.Number ", p.Number)
		// fmt.Print(" p.UserName ", p.UserName)
		// fmt.Print(" p.Time ", p.Time)
		// fmt.Println(" p.Content ", p.Content)

	})

	return posts
}

func (t *Thread) fetchThread() *scraper.Scraper {

	s := scraper.NewScraper("https://vozforums.com/showthread.php?t="+strconv.Itoa(t.ID), "utf-8")
	return s
}

// Posts is the list of Posts in the Thread
type Posts []Post
