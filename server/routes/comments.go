package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbHelper "github.com/treavorgagne/twitter-clone/server/db"
	"github.com/treavorgagne/twitter-clone/server/models"
)


func CreateComment(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	// Bind incoming JSON to the struct
	var req models.CreateCommentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}	
	user_id := c.Param("user_id")
	tweet_id := c.Param("tweet_id")

	res, err := conn.ExecContext(c.Request.Context(), "insert into comments (body, tweet_id, user_id) values (?, ?, ?);", req.Body, tweet_id, user_id)
	if err != nil {
		log.Panic("ERROR_CREATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Fetch the new user ID
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("ERROR_FETCHING_INSERT_ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("users %s comment on tweet %s", user_id, tweet_id)
	c.JSON(http.StatusOK, gin.H{"comment_id": id})
}

func UpdateComment(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	// Bind incoming JSON to the struct
	var req models.UpdateCommentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}	

	comment_id := c.Param("comment_id")
	_, err = conn.ExecContext(c.Request.Context(), "UPDATE comments SET body = ? WHERE comment_id = ? ;", req.Body, comment_id)
	if err != nil {
		log.Panic("ERROR_CREATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.Status(http.StatusNoContent) 
}

func DeleteComment(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	comment_id := c.Param("comment_id")
	_, err = conn.ExecContext(c.Request.Context(), "delete from comments where comment_id = ?;", comment_id)
	if err != nil {
		log.Panic("ERROR_CREATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func LikeComment(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	comment_id := c.Param("comment_id")
	user_id := c.Param("user_id")
	_, err = conn.ExecContext(c.Request.Context(), "insert into likes_comment (comment_id, user_id) values (?, ?);", comment_id, user_id)
	if err != nil {
		log.Panic("ERROR_CREATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func UnLikeComment(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	comment_id := c.Param("comment_id")
	user_id := c.Param("user_id")
	_, err = conn.ExecContext(c.Request.Context(), "delete from likes_comment where comment_id = ? and user_id = ?;", comment_id, user_id)
	if err != nil {
		log.Panic("ERROR_CREATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.Status(http.StatusNoContent) 
}

