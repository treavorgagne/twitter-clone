package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbHelper "github.com/treavorgagne/twitter-clone/server/db"
)

func FollowUser(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("user_id")
	follow_id := c.Param("follow_id")

	_, err = conn.ExecContext(c.Request.Context(), "insert into follows (user_id, follow_id) values (?, ?);", user_id, follow_id)
	if err != nil {
		log.Panic("ERROR_FOLLOWING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error trying to follow user"})
		return
	}

	log.Printf(`user_id: %s followed user_id: %s`, user_id, follow_id)
	c.Status(http.StatusCreated)
}

func UnFollowUser(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("user_id")
	follow_id := c.Param("follow_id")

	_, err = conn.ExecContext(c.Request.Context(),"delete from follows where user_id = ? and follow_id = ?;", user_id, follow_id)
	if err != nil {
		log.Panic("ERROR_UNFOLLOWING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error trying to unfollow user"})
		return
	}

	log.Printf(`user_id: %s unfollowed user_id: %s`, user_id, follow_id)
	c.Status(http.StatusNoContent)
}