package vozscrape

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nnti3n/voz-archive-plus/scraper"
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
func NewThread(id string, title string, source string, pageCount string) *Thread {

	t := new(Thread)
	t.ID, _ = strconv.Atoi(strings.TrimSpace(id))
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

	pPage.Find("#posts > div").Each(func(i int, s *goquery.Selection) {
		p := new(Post)
		number := s.Find("tr:first-child td div:first-child a:first-child")
		p.Number, _ = strconv.Atoi(number.Text())
		_number, _ := number.Attr("href")
		p.PostID, _ = strconv.Atoi(strings.Split(strings.Split(_number, "=")[1], "&")[0])

		p.Time = s.Find("tr:first-child td.thead div:nth-child(2)").Text()

		username, _ := s.Find(".bigusername").Attr("href")
		p.UserName = strings.Split(username, "u=")[1]

		p.Content = s.Find(".voz-post-message").Text()

	})

	return posts
}

func (t *Thread) fetchThread() *scraper.Scraper {

	s := scraper.NewScraper("https://vozforums.com/showthread?t="+strconv.Itoa(t.ID), "utf-8")
	return s
}

// Posts is the list of Posts in the Thread
type Posts []Post
