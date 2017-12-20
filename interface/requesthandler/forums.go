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
	threadID, _ := strconv.Atoi(c.Param("threadID"))
	log.Println(c.Param("threadID"))

	thread := vozscrape.Thread{ID: threadID}
	err := e.Db.Select(&thread)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"querydata": thread,
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
