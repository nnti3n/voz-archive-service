package requesthandler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/nnti3n/voz-archive-plus/serviceWorker/vozscrape"

	"github.com/gin-gonic/gin"
)

// Env store db
type Env struct {
	Db *pg.DB
}

// FetchAllThread fetch all threads
func (e *Env) FetchAllThread(c *gin.Context) {
	boxID := c.Param("boxID")

	threads := []vozscrape.Thread{}
	err := e.Db.Model(&threads).Where("box_id = ?", boxID).Limit(20).Select()
	log.Println(c.Param("boxID"))
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

// FetchSingleThread fetch all posts of thread
func (e *Env) FetchSingleThread(c *gin.Context) {
	threadID, err := strconv.Atoi(c.Param("ThreadID"))

	thread := vozscrape.Thread{ID: 6778933}
	err = e.Db.Select(&thread)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"querydata": threadID,
			"data":      thread,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"err":       err,
			"querydata": threadID,
			"data":      []string{},
		})
	}
}
