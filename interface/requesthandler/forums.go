package requesthandler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/nnti3n/voz-archive-plus/serviceWorker/vozscrape"

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

// threadFilter filter threads by id and page
func (f *ThreadFilter) threadFilter(q *orm.Query) (*orm.Query, error) {
	if f.BoxID > 0 {
		q = q.Where("box_id = ?", f.BoxID)
	}
	q = q.Apply(f.Pager.Paginate)
	return q, nil
}

// FetchAllThread fetch all threads
func (e *Env) FetchAllThread(c *gin.Context) {
	var filter ThreadFilter
	filter.Pager.SetURLValues(c.Request.URL.Query())
	filter.BoxID, _ = strconv.Atoi(c.Param("boxID"))

	threads := []vozscrape.Thread{}
	err := e.Db.Model(&threads).Apply(filter.threadFilter).Select()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":   threads,
			"params": filter.BoxID,
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
	var filter PostFilter
	filter.Pager.SetURLValues(c.Request.URL.Query())
	filter.ThreadID, _ = strconv.Atoi(c.Param("threadID"))

	posts := []vozscrape.Post{}
	err := e.Db.Model(&posts).Apply(filter.postFilter).Select()

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":   posts,
			"params": filter.ThreadID,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"err":  err,
			"data": []string{},
		})
	}
}
