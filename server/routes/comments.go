package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	dbHelper "github.com/treavorgagne/twitter-clone/server/db"
	"github.com/treavorgagne/twitter-clone/server/models"
)

func GetComment(c *gin.Context, db *sql.DB, rdb *redis.Client) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	comment_id := c.Param("comment_id")
	var comment models.GetCommentsResponse

	row := conn.QueryRowContext(c.Request.Context(), "select * from comments_stats where comment_id = ?;", comment_id)
	err = row.Scan(&comment.Comment_id, &comment.Tweet_id, &comment.Body, &comment.User_id, &comment.Created_At, &comment.Comment_total_likes);
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	} else if err != nil {
		log.Printf("ERROR_SCANNING_TWEET_ROW: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tweet data"})
		return
	}

	// 🔥 Store into Redis cache
	// 🧹 Marshal struct to JSON before caching
	commentJSON, err := json.Marshal(comment)
	if err != nil {
		log.Println("Failed to marshal user:", err)
	} else {
		cacheKey := c.Request.URL.Path
		err := rdb.SetNX(c, cacheKey, commentJSON, 3*time.Minute).Err()
		if err != nil {
			log.Println("Redis cache SET error:", err)
		}
	}

	c.JSON(http.StatusOK, comment)
}

func GetCommentsByTweetId(c *gin.Context, db *sql.DB, rdb *redis.Client) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("tweet_id")

	rows, err := conn.QueryContext(c.Request.Context(), "select * from comments_stats where tweet_id = ?;", user_id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Comments found"})
		return
	} else if err != nil {
		log.Panic("ERROR_GETTING_COMMENTS_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	defer rows.Close()
	var comments []models.GetCommentsResponse
	for rows.Next() {
		var comment models.GetCommentsResponse
		if err = rows.Scan(&comment.Comment_id, &comment.Tweet_id, &comment.Body, &comment.User_id, &comment.Created_At, &comment.Comment_total_likes); err != nil {
			log.Printf("ERROR_SCANNING_TWEET_ROW: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tweet data"})
			return
		}
		comments = append(comments, comment)
	}

	// 🔥 Store into Redis cache
	// 🧹 Marshal struct to JSON before caching
	commentsJSON, err := json.Marshal(comments)
	if err != nil {
		log.Println("Failed to marshal user:", err)
	} else {
		cacheKey := c.Request.URL.Path
		err := rdb.SetNX(c, cacheKey, commentsJSON, 3*time.Minute).Err()
		if err != nil {
			log.Println("Redis cache SET error:", err)
		}
	}

	c.JSON(http.StatusOK, comments)
}

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

