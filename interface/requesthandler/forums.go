package requesthandler

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/nnti3n/voz-archive-service/serviceWorker/vozscrape"

	"github.com/gin-gonic/gin"
)

// Env store db
type Env struct {
	Db *pg.DB
}

// ThreadFilter struct for filter
type ThreadFilter struct {
	orm.Pager
	BoxID int
}

var excludeThreads = []int{6735473, 6609261, 3613304, 2024506, 5523490}

// threadFilter filter threads by id and page
func (f *ThreadFilter) threadFilter(q *orm.Query) (*orm.Query, error) {
	if f.BoxID > 0 {
		q = q.Where("box_id = ?", f.BoxID).
			Where("id not in (?)", pg.In(excludeThreads)).
			Where("post_count > ?", 0)
	}
	q = q.Apply(f.Pager.Paginate)
	return q, nil
}

// FetchAllThread fetch all threads
func (e *Env) FetchAllThread(c *gin.Context) {
	var filter ThreadFilter
	filter.Pager.SetURLValues(c.Request.URL.Query())
	boxID, _ := strconv.Atoi(c.Param("boxID"))
	filter.BoxID = boxID

	threads := []vozscrape.Thread{}
	count, err := e.Db.Model(&threads).Apply(filter.threadFilter).
		SelectAndCount()
	pageCount := count / 10
	if math.Mod(float64(count), 10) > 0 {
		pageCount = count/10 + 1
	}
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":   threads,
			"params": filter.BoxID,
			"page":   pageCount,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"data":   []string{},
			"params": filter.BoxID,
		})
	}

}

// FetchSingleThread fetch info of single thread
func (e *Env) FetchSingleThread(c *gin.Context) {
	threadID, _ := strconv.Atoi(c.Param("threadID"))
	log.Println(c.Param("threadID"))

	thread := vozscrape.Thread{ID: threadID}
	err := e.Db.Select(&thread)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": thread,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"err":  err,
			"data": []string{},
		})
	}
}

// PostFilter struct for filter
type PostFilter struct {
	orm.Pager
	ThreadID int
}

// postFilter filter posts by id and page
func (f *PostFilter) postFilter(q *orm.Query) (*orm.Query, error) {
	if f.ThreadID > 0 {
		q = q.Where("thread_id = ?", f.ThreadID)
	}
	q = q.Order("number ASC").Apply(f.Pager.Paginate)
	return q, nil
}

// FetchThreadPosts fetch all posts of thread
func (e *Env) FetchThreadPosts(c *gin.Context) {
	threadID, _ := strconv.Atoi(c.Param("threadID"))
	var filter PostFilter
	filter.Pager.SetURLValues(c.Request.URL.Query())
	filter.ThreadID = threadID

	posts := []vozscrape.Post{}
	count, err := e.Db.Model(&posts).Apply(filter.postFilter).
		SelectAndCount()
	pageCount := count / 10
	if math.Mod(float64(count), 10) > 0 {
		pageCount = count/10 + 1
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":   posts,
			"params": filter.ThreadID,
			"page":   pageCount,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"err":  err,
			"data": []string{},
		})
	}
}
