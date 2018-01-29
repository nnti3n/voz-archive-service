package requesthandler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/nnti3n/voz-archive-plus/serviceWorker/vozscrape"
	"github.com/nnti3n/voz-archive-plus/utilities"

	"github.com/gin-gonic/gin"
)

// Env store db
type Env struct {
	Db *pg.DB
}

// FetchAllThread fetch all threads
func (e *Env) FetchAllThread(c *gin.Context) {
	boxID := c.Param("boxID")
	limit, offset := utilities.Pagination(c, 20)

	threads := []vozscrape.Thread{}
	err := e.Db.Model(&threads).Where("box_id = ?", boxID).
		Offset(offset).Limit(limit).Order("id ASC").Select()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":   threads,
			"params": boxID,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"data":   []string{},
			"params": boxID,
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
			"err":    err,
			"params": threadID,
			"data":   []string{},
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
	// limit, offset := utilities.Pagination(c, 20)

	var filter PostFilter
	filter.Pager.SetURLValues(c.Request.URL.Query())
	filter.ThreadID = threadID

	posts := []vozscrape.Post{}
	// err := e.Db.Model(&posts).Where("thread_id = ?", threadID).
	// 	Offset(offset).Limit(limit).Order("number ASC").Select()
	err := e.Db.Model(&posts).Apply(filter.postFilter).Select()

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":   posts,
			"params": threadID,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"err":    err,
			"params": threadID,
			"data":   []string{},
		})
	}
}
