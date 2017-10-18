package vozscrape

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/ref/mercato/scraper"
)

// Thread is the model for the forums threads
type Thread struct {
	url         string
	title       string
	id          int
	source      string
	pageCount   int
	lastUpdated string

	Posts []*Post `json:"results"`
}

// Post is the struct for a single Post in the Grocery Store
type Post struct {
	url      string
	PostID   int    `json:"postid"`
	UserID   int    `json:"userid"`
	UserName string `json:"username"`
	Time     string `json:"time"`
	Content  string `json:"content"`
}

// NewPost creates a Post and fills missing information
// from the Post page
func NewPost(url string, postID string, userID string, time string, content string) *Post {

	p := new(Post)
	p.url = strings.TrimSpace(url)
	pPage := p.fetchPage()
	p.PostID, _ = strconv.Atoi(strings.TrimSpace(postID))
	p.PostID, _ = strconv.Atoi(strings.TrimSpace(userID))

	p.Description = p.getDescription(pPage)

	return p
}

func (p *Post) getDescription(pPage *scraper.Scraper) string {
	d := ""

	pPage.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "description" {
			description, _ := s.Attr("content")
			d = description
		}
	})

	return strings.TrimSpace(d)
}

func (p *Post) fetchPage() *scraper.Scraper {

	s := scraper.NewScraper(p.url, "utf-8")
	return s
}

// Posts is the list of Posts in the Grocery Store
type Posts []Post
