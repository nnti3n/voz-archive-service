package vozscrape

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nnti3n/voz-archive-plus/scraper"
)

// Thread is the model for the forums threads
type Thread struct {
	ID        int
	Title     string
	Source    string
	PageCount int
	PostCount int
	ViewCount int
	BoxID     int

	Posts []*Post
}

// Post is the struct for a single Post in thread
type Post struct {
	ID       int
	Number   int
	UserID   int
	UserName string
	Time     string
	Content  string
	ThreadID int
}

// NewThread creates a Thread and fills missing information
// from the Thread page
func NewThread(id int, title string, source string, pageCount string, postCount string, viewCount string, boxID int) *Thread {

	t := new(Thread)
	t.ID = id
	t.Title = title
	t.Source = source
	t.PageCount, _ = strconv.Atoi(strings.Replace(pageCount, ",", "", -1))
	t.PostCount, _ = strconv.Atoi(strings.Replace(postCount, ",", "", -1))
	t.ViewCount, _ = strconv.Atoi(strings.Replace(viewCount, ",", "", -1))
	t.BoxID = boxID

	// Start scraping thread
	tPage := t.fetchThread()

	t.Posts = t.getPosts(tPage)

	return t
}

func (t *Thread) getPosts(pPage []*scraper.Scraper) []*Post {
	posts := []*Post{}

	for _, page := range pPage {
		page.Find("#posts > div [align='left']").Each(func(i int, s *goquery.Selection) {
			p := new(Post)
			number := s.Find("tr:first-child td div:first-child a:first-child")
			numberName, exist := number.Attr("name")
			if !exist {
				fmt.Println("Not found page number, set -1")
				numberName = "-1"
			}
			p.Number, _ = strconv.Atoi(numberName)
			_number, exist := number.Attr("href")
			if !exist {
				fmt.Println("no post href")
				return
			}
			p.ID, _ = strconv.Atoi(strings.Split(strings.Split(_number, "=")[1], "&")[0])

			p.Time = strings.TrimSpace(s.Find("tr:first-child td.thead div:nth-child(2)").Text())

			userID, exist := s.Find(".bigusername").Attr("href")
			if !exist {
				fmt.Println("not found userID")
			}
			p.UserID, _ = strconv.Atoi(strings.Split(userID, "u=")[1])
			p.UserName = strings.TrimSpace(s.Find(".bigusername").Text())
			p.Content = strings.TrimSpace(s.Find(".voz-post-message").Text())
			p.ThreadID = t.ID

			posts = append(posts, p)
		})
	}

	return posts
}

func (t *Thread) fetchThread() []*scraper.Scraper {
	s := []*scraper.Scraper{}
	for i := 1; i < t.PageCount; i++ {
		p := scraper.NewScraper("https://vozforums.com/showthread.php?t="+strconv.Itoa(t.ID)+"&page="+strconv.Itoa(t.PageCount), "utf-8")
		s = append(s, p)
	}
	return s
}

// Posts is the list of Posts in the Thread
type Posts []Post
