package vozscrape

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nnti3n/voz-archive-service/serviceWorker/scraper"
	"github.com/nnti3n/voz-archive-service/utilities"
)

// Thread is the model for the forums threads
type Thread struct {
	ID              int
	Title           string
	Source          string
	PageCount       int
	PostCount       int
	ViewCount       int
	BoxID           int
	UserIDStarter   int
	UserNameStarter string

	Posts []*Post `sql:"-"`
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

var excludeThreads = []int{6735473, 6609261, 3613304, 2024506, 5523490}

// NewThread creates a Thread and fills missing information
// from the Thread page
func NewThread(id int, title string, userID int, userName string, source string, pageCount string, postCount string, viewCount string, boxID int) *Thread {

	t := new(Thread)
	t.ID = id
	t.Title = title
	t.UserIDStarter = userID
	t.UserNameStarter = userName
	t.Source = source
	t.PageCount, _ = strconv.Atoi(strings.Replace(pageCount, ",", "", -1))
	t.PostCount, _ = strconv.Atoi(strings.Replace(postCount, ",", "", -1))
	t.ViewCount, _ = strconv.Atoi(strings.Replace(viewCount, ",", "", -1))
	t.BoxID = boxID

	if utilities.NumberInSlice(id, excludeThreads) {
		log.Println("met", id)
		return t
	}

	thread := Thread{}
	err := db.Model(&thread).Where("id = ?", id).Select()
	if err != nil {
		thread.PageCount = 1
	}
	if thread.PostCount == t.PostCount {
		log.Println("same postcount", id)
		return t
	}

	// Start scraping thread
	tPage := t.fetchThread(thread.PageCount)

	t.Posts = t.getPosts(tPage)

	return t
}

func (t *Thread) getPosts(pPage []scraper.Scraper) []*Post {
	posts := []*Post{}

	for _, page := range pPage {
		page.Find("#posts > div [align='left']").Each(func(i int, s *goquery.Selection) {
			p := new(Post)
			number := s.Find("tr:first-child td div:first-child a:first-child")
			numberName, exist := number.Attr("name")
			if !exist {
				fmt.Println("Not found post number, set -1")
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
				p.UserID = -1
			} else {
				p.UserID, _ = strconv.Atoi(strings.Split(userID, "u=")[1])
			}
			p.UserName = strings.TrimSpace(s.Find(".bigusername").Text())
			_content, _ := s.Find(".voz-post-message").Html()
			p.Content = strings.TrimSpace(_content)
			p.ThreadID = t.ID

			posts = append(posts, p)
		})
	}

	return posts
}

func (t *Thread) fetchThread(currentPageCount int) []scraper.Scraper {
	s := []scraper.Scraper{}
	i := 1
	if t.PageCount >= currentPageCount {
		i = currentPageCount
		log.Println("Count", currentPageCount, t.PageCount)
	}
	for i <= t.PageCount {
		p := scraper.NewScraper("https://vozforums.com/showthread.php?t="+strconv.Itoa(t.ID)+"&page="+strconv.Itoa(i), "utf-8")
		s = append(s, *p)
		i++
	}
	return s
}

// Posts is the list of Posts in the Thread
type Posts []Post
