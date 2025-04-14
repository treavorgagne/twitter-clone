package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FollowUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		follow_id := c.Param("follow_id")

		_, err := db.Exec("insert into follows (user_id, follow_id) values (?, ?);", user_id, follow_id)
		if err != nil {
			log.Panic("ERROR_FOLLOWING_USER_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error trying to follow user"})
			return
		}

		log.Printf(`user_id: %s followed user_id: %s`, user_id, follow_id)
		c.Status(http.StatusCreated)
	}
}

func UnFollowUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		follow_id := c.Param("follow_id")

		_, err := db.Exec("delete from follows where user_id = ? and follow_id = ?;", user_id, follow_id)
		if err != nil {
			log.Panic("ERROR_UNFOLLOWING_USER_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error trying to unfollow user"})
			return
		}

		log.Printf(`user_id: %s unfollowed user_id: %s`, user_id, follow_id)
		c.Status(http.StatusNoContent)
	}
}