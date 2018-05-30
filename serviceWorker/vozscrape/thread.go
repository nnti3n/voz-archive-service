package vozscrape

import (
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-pg/pg"
	"github.com/nnti3n/voz-archive-service/serviceWorker/scraper"
	"github.com/nnti3n/voz-archive-service/utilities"
)

// Thread is the model for the forums threads
type Thread struct {
	ID              int
	Title           string
	Source          string
	PageCount       int `sql:",notnull"`
	PostCount       int `sql:",notnull"`
	ViewCount       int `sql:",notnull"`
	BoxID           int
	UserIDStarter   int
	UserNameStarter string
	LastUpdated     time.Time `sql:",notnull"`

	Posts []*Post `sql:"-"`
}

// Post is the struct for a single Post in thread
type Post struct {
	ID       int
	Number   int `sql:",notnull"`
	UserID   int
	UserName string
	Time     time.Time `sql:",notnull"`
	Content  string
	ThreadID int
}

var excludeThreads = []int{6735473, 6609261, 3613304, 2024506, 5523490, 4289698,
	3805350, 3805345, 3973882, 4000407, 2926112, 6830409, 4573733, 6822821, 4080256}

// NewThread creates a Thread and fills missing information
// from the Thread page
func NewThread(id int, title string, userID int, userName string, source string,
	pageCount string, postCount string, viewCount string, boxID int) *Thread {

	t := new(Thread)
	t.ID = id
	t.Title = title
	t.UserIDStarter = userID
	t.UserNameStarter = userName
	t.LastUpdated = time.Now()
	t.Source = source
	t.PageCount, _ = strconv.Atoi(strings.Replace(pageCount, ",", "", -1))
	t.PostCount, _ = strconv.Atoi(strings.Replace(postCount, ",", "", -1))
	t.ViewCount, _ = strconv.Atoi(strings.Replace(viewCount, ",", "", -1))
	t.BoxID = boxID

	if utilities.NumberInSlice(id, excludeThreads) {
		// log.Println("met", id)
		return t
	}

	thread := Thread{}
	err := db.Model(&thread).Where("id = ?", id).Select()

	if err != nil {
		thread.PageCount = 1
	}

	var count int
	_, errCount := db.Model((*Post)(nil)).
		QueryOne(pg.Scan(&count), `SELECT count(*) FROM posts WHERE thread_id = ?`, id)
	if errCount != nil {
		count = 0
	}

	// only scrape max 20 page
	if t.PageCount >= count/10+20 {
		t.PageCount = count/10 + 20
		t.PostCount = count + 20*10
	}

	if count == t.PostCount {
		// log.Println("same postcount", thread.ID, thread.PostCount)
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
		page.Find("#posts > div [align='left']").
			Each(func(i int, s *goquery.Selection) {
				p := new(Post)
				number := s.Find("tr:first-child td div:first-child a:first-child")
				numberName, exist := number.Attr("name")
				if !exist {
					numberName = "-1"
				}
				p.Number, _ = strconv.Atoi(numberName)
				_number, exist := number.Attr("href")
				if !exist {
					return
				}
				p.ID, _ = strconv.
					Atoi(strings.Split(strings.Split(_number, "=")[1], "&")[0])
				p.Time = utilities.ParseTime(strings.
					TrimSpace(s.Find("tr:first-child td.thead div:nth-child(2)").
						Text()))

				userID, exist := s.Find(".bigusername").Attr("href")
				if !exist {
					p.UserID = -1
				} else {
					p.UserID, _ = strconv.Atoi(strings.Split(userID, "u=")[1])
				}
				p.UserName = strings.TrimSpace(s.Find(".bigusername").Text())

				// dealing with post content
				_content, _ := s.Find(".voz-post-message").Html()
				rcontent := regexp.MustCompile(`"(\/redirect\/index.php)(.*?)"`)
				urlText := rcontent.FindString(_content)
				if urlText != "" {
					urlReplace, _ := url.QueryUnescape(urlText)
					urlReplace = strings.Replace(urlReplace, "/redirect/index.php?link=", "", -1)
					_content = strings.Replace(_content, urlText, urlReplace, -1)
				}
				p.Content = strings.TrimSpace(_content)

				p.ThreadID = t.ID
				posts = append(posts, p)
			})

		// end of page
	}

	log.Println(t.ID, len(posts))
	return posts
}

func (t *Thread) fetchThread(currentPageCount int) []scraper.Scraper {
	s := []scraper.Scraper{}

	// scrape
	count := currentPageCount
	for count <= t.PageCount {
		p := scraper.
			NewScraper("https://vozforums.com/showthread.php?t="+
				strconv.Itoa(t.ID)+"&page="+strconv.Itoa(count), "utf-8")
		s = append(s, *p)
		count++
	}
	return s
}

// Posts is the list of Posts in the Thread
type Posts []Post
